/*
------------------------------------------------------------------------------------------------------------------------
####### Copyright (c) 2022-2023 Archivage Numérique.
####### All rights reserved.
####### Use of this source code is governed by a MIT style license that can be found in the LICENSE file.
------------------------------------------------------------------------------------------------------------------------
*/

package dikit

type (
	closableFactory interface {
		close() error
	}

	closable interface {
		Close() error
	}
)

func empty[T any]() (t T) {
	return
}

/*
######################################################################################################## @(^_^)@ #######
*/
