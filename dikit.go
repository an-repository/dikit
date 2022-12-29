/*
------------------------------------------------------------------------------------------------------------------------
####### Copyright (c) 2022-2023 Archivage Numérique.
####### All rights reserved.
####### Use of this source code is governed by a MIT style license that can be found in the LICENSE file.
------------------------------------------------------------------------------------------------------------------------
*/

package dikit

import "github.com/an-repository/errors"

func Add[T any](c *Container, name string, b Builder[T]) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.add(name, newFactoryWithBuilder(name, b))
}

func AddValue[T any](c *Container, name string, value T) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.add(name, newFactoryWithValue(name, value)); err != nil {
		return err
	}

	c.addToClose(name)

	return nil
}

func Get[T any](c *Container, name string) (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	f, ok := c.get(name)
	if !ok {
		return empty[T](), errors.New("factory not found", "name", name) ///////////////////////////////////////////////
	}

	factory, ok := f.(*factory[T])
	if !ok {
		return empty[T](), errors.New("non-conforming factory", "name", name) //////////////////////////////////////////
	}

	instance, err := factory.getInstance(c)
	if err != nil {
		return empty[T](), err
	}

	return instance, nil
}

/*
######################################################################################################## @(^_^)@ #######
*/
