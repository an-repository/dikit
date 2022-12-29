/*
------------------------------------------------------------------------------------------------------------------------
####### Copyright (c) 2022-2023 Archivage Numérique.
####### All rights reserved.
####### Use of this source code is governed by a MIT style license that can be found in the LICENSE file.
------------------------------------------------------------------------------------------------------------------------
*/

package dikit

import (
	"strings"
	"sync"

	"github.com/an-repository/errors"
	"github.com/an-repository/tracer"
)

type Container struct {
	factories map[string]any
	toClose   []string
	mutex     sync.Mutex
}

func NewContainer() *Container {
	return &Container{
		factories: make(map[string]any),
		toClose:   make([]string, 0),
	}
}

func (c *Container) add(name string, f any) error {
	_, ok := c.factories[name]
	if ok {
		return errors.New("this factory already exists", "name", name) /////////////////////////////////////////////////
	}

	tracer.Sendf("[dikit] add(%s)", name) //............................................................................

	c.factories[name] = f

	return nil
}

func (c *Container) get(name string) (any, bool) {
	f, ok := c.factories[name]
	return f, ok
}

func (c *Container) addToClose(name string) {
	c.toClose = append(c.toClose, name)
}

func (c *Container) close(name string) error {
	return c.factories[name].(closableFactory).close()
}

func (c *Container) Close() error {
	merr := make([]string, 0)

	for i := len(c.toClose) - 1; i >= 0; i-- {
		if err := c.close(c.toClose[i]); err != nil {
			merr = append(merr, err.Error())
		}
	}

	if len(merr) == 0 {
		return nil
	}

	return errors.New(strings.Join(merr, " ### ")) /////////////////////////////////////////////////////////////////////
}

/*
######################################################################################################## @(^_^)@ #######
*/
