// Copyright (c) 2014 Yeung Shu Hung (Koala Yeung)
//
//  This file is part of RESTit.
//
//  RESTit is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  RESTit is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  Use of this source code is governed by the GPL v3 license. A copy
//  of the licence can be found in the LICENSE.md file along with RESTit.
//  If not, see <http://www.gnu.org/licenses/>.

package restit

import (
	"fmt"
	"github.com/jmcvetta/napping"
)

type Case struct {
	Request      *napping.Request
	Session      Session
	Expectations []Expectation
	Tester       *Tester
}

// To actually run the test case
func (c *Case) Run() (r *Result, err error) {
	res, err := c.Session.Send(c.Request)
	result := Result{
		Response: res,
	}
	r = &result

	// test each expectations
	resp := (*c).Request.Result.(Response)
	for i := 0; i < len(c.Expectations); i++ {
		err = c.Expectations[i].Test(resp)
		if err != nil {
			err = fmt.Errorf("Failed in test: \"%s\" "+
				"Reason: \"%s\"",
				c.Expectations[i].Desc,
				err.Error())
			return
		}
	}

	return
}

// To run the test case and panic on error
func (c *Case) RunOrPanic() (r *Result) {
	r, err := c.Run()
	if err != nil {
		panic(err)
	}
	return
}

// Set the result to the given interface{}
func (c *Case) WithResponseAs(r interface{}) *Case {
	c.Request.Result = r
	return c
}

// Set the query parameter
func (c *Case) WithParams(p *napping.Params) *Case {
	c.Request.Params = p
	return c
}

// Append Test to Expectations
// Tests if the result count equal to n
func (c *Case) ExpectResultCount(n int) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Test Result Count",
		Test: func(r Response) (err error) {
			count := r.Count()
			if count != n {
				err = fmt.Errorf(
					"Result count is %d "+
						"(expected %d)",
					count, n)
			}
			return
		},
	})
	return c
}

// Append Test to Expectations
// Tests if the result count not equal to n
func (c *Case) ExpectResultCountNot(n int) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Test Result Count",
		Test: func(r Response) (err error) {
			count := r.Count()
			if count == n {
				err = fmt.Errorf(
					"Result count is %d "+
						"(expected %d)",
					count, n)
			}
			return
		},
	})
	return c
}

// Append Test to Expectations
// Tests if the item is valid
func (c *Case) ExpectResultsValid() *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: "Test Results Valid",
		Test: func(r Response) (err error) {
			for i := 0; i < r.Count(); i++ {
				err = r.NthValid(i)
				if err != nil {
					err = fmt.Errorf(
						"Item %d invalid: %s",
						i, err.Error())
					return
				}
			}
			return
		},
	})
	return c
}

// Append Test to Expectation
// Tests if the nth item matches the provided one
func (c *Case) ExpectResultNth(n int, b interface{}) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: fmt.Sprintf("Test #%d Result Valid", n),
		Test: func(r Response) (err error) {
			a, err := r.GetNth(n)
			if err != nil {
				return
			}
			err = r.Match(a, b)
			return
		},
	})
	return c
}

// Append Custom Test to Expectation
// Allow user to inject user defined tests
func (c *Case) ExpectResultsToPass(
	desc string, test func(Response) error) *Case {
	c.Expectations = append(c.Expectations, Expectation{
		Desc: desc,
		Test: test,
	})
	return c
}

// Expection to the response in a Case
type Expectation struct {
	Desc string
	Test func(Response) error
}

// Test Result of a Case
type Result struct {
	Response *napping.Response
}

// Wrap the napping.Session in Session
// to make unit testing easier
type Session interface {
	Send(*napping.Request) (*napping.Response, error)
}
