/*
------------------------------------------------------------------------------------------------------------------------
####### Copyright (c) 2022-2023 Archivage Numérique.
####### All rights reserved.
####### Use of this source code is governed by a MIT style license that can be found in the LICENSE file.
------------------------------------------------------------------------------------------------------------------------
*/

package dikit

import (
	"runtime/debug"
	"sync"

	"github.com/an-repository/errors"
	"github.com/an-repository/tracer"
)

type factory[T any] struct {
	name     string
	instance T
	built    bool
	builder  Builder[T]
	mutex    sync.Mutex
}

func newFactoryWithValue[T any](name string, value T) *factory[T] {
	return &factory[T]{
		name:     name,
		instance: value,
		built:    true,
	}
}

func newFactoryWithBuilder[T any](name string, b Builder[T]) *factory[T] {
	return &factory[T]{
		name:    name,
		builder: b,
	}
}

func (f *factory[T]) build(c *Container) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New( //////////////////////////////////////////////////////////////////////////////////////
					"recovered error when building the instance",
					"factory", f.name,
					"data", r,
					"stack", debug.Stack(),
				)
			}
		}
	}()

	tracer.Sendf("[dikit] build(%s)", f.name) //........................................................................

	instance, err := f.builder(c)
	if err == nil {
		f.instance = instance
		f.built = true

		c.addToClose(f.name)

		return nil
	}

	return errors.WithMessage(err, "unable to build", "factory", f.name) ///////////////////////////////////////////////
}

func (f *factory[T]) getInstance(c *Container) (T, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !f.built {
		if err := f.build(c); err != nil {
			return empty[T](), err
		}
	}

	return f.instance, nil
}

func (f *factory[T]) close() error {
	instance, ok := any(f.instance).(closable)
	if !ok {
		return nil
	}

	tracer.Sendf("[dikit] close(%s)", f.name) //........................................................................

	return instance.Close()
}

/*
######################################################################################################## @(^_^)@ #######
*/
