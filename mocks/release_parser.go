// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	colly "github.com/gocolly/colly/v2"
	crawler "github.com/hostfactor/minecrafter/crawler"
	mock "github.com/stretchr/testify/mock"
)

// ReleaseParser is an autogenerated mock type for the ReleaseParser type
type ReleaseParser struct {
	mock.Mock
}

// ParseRelease provides a mock function with given fields: e
func (_m *ReleaseParser) ParseRelease(e *colly.HTMLElement) *crawler.Release {
	ret := _m.Called(e)

	var r0 *crawler.Release
	if rf, ok := ret.Get(0).(func(*colly.HTMLElement) *crawler.Release); ok {
		r0 = rf(e)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*crawler.Release)
		}
	}

	return r0
}
