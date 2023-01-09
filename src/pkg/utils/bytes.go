// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors
// forked from https://www.socketloop.com/tutorials/golang-byte-format-example

// Package utils provides generic helper functions.
package utils

import (
	"math"
	"strconv"
)

// RoundUp rounds a float64 to the given number of decimal places.
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

// ByteFormat formats a number of bytes into a human readable string.
func ByteFormat(inputNum float64, precision int) string {
	if precision <= 0 {
		precision = 1
	}

	var unit string
	var returnVal float64

	// https://www.techtarget.com/searchstorage/definition/mebibyte-MiB
	if inputNum >= 1000000000 {
		returnVal = RoundUp(inputNum/1000000000, precision)
		unit = " GB" // gigabyte
	} else if inputNum >= 1000000 {
		returnVal = RoundUp(inputNum/1000000, precision)
		unit = " MB" // megabyte
	} else if inputNum >= 1000 {
		returnVal = RoundUp(inputNum/1000, precision)
		unit = " KB" // kilobyte
	} else {
		returnVal = inputNum
		unit = " Byte" // byte
	}

	if returnVal > 1 {
		unit += "s"
	}

	return strconv.FormatFloat(returnVal, 'f', precision, 64) + unit
}
