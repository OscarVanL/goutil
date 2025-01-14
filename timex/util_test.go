package timex_test

import (
	"testing"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/timex"
)

func TestBasic(t *testing.T) {
	sec := timex.NowUnix()

	assert.NotEmpty(t, timex.FormatUnix(sec))
	assert.NotEmpty(t, timex.FormatUnixBy(sec, time.RFC3339))

	tt := timex.TodayStart()
	assert.Eq(t, "00:00:00", timex.DateFormat(tt, "H:I:S"))

	tt = timex.TodayEnd()
	assert.Eq(t, "23:59:59", timex.DateFormat(tt, "H:I:S"))

	tt = timex.NowHourStart()
	assert.Eq(t, "00:00", timex.DateFormat(tt, "I:S"))

	tt = timex.NowHourEnd()
	assert.Eq(t, "59:59", timex.DateFormat(tt, "I:S"))
}

func TestDateFormat(t *testing.T) {
	now := time.Now()

	tests := []struct{ layout, template string }{
		{"20060102 15:04:05", "Ymd H:I:S"},
		{"2006-01-02 15:04:05", "Y-m-d H:I:S"},
		{"2006-01-02 15:04", "Y-m-d H:I"},
		{"01/02 15:04:05", "m/d H:I:S"},
		{"06/01/02 15:04:05", "y/m/d H:I:S"},
		{"06/01/02 15:04:05.000", "y/m/d H:I:Sv"},
	}

	for i, item := range tests {
		date := now.Format(item.layout)
		assert.Eq(t, date, timex.DateFormat(now, item.template))
		if i%2 == 0 {
			assert.Eq(t, date, timex.Date(now, item.template))
		}
	}

	assert.Eq(t, now.Format("01/02 15:04:05.000000"), timex.Date(now, "m/d H:I:Su"))
}

func TestFormatUnix(t *testing.T) {
	now := time.Now()
	want := now.Format("2006-01-02 15:04:05")

	assert.Eq(t, want, timex.FormatUnix(now.Unix()))
	assert.Eq(t, want, timex.FormatUnixBy(now.Unix(), timex.DefaultLayout))
	assert.Eq(t, want, timex.FormatUnixByTpl(now.Unix(), "Y-m-d H:I:S"))
	dump.P(want)
}

func TestToLayout(t *testing.T) {
	assert.Eq(t, timex.DefaultLayout, timex.ToLayout(""))
	assert.Eq(t, time.RFC3339, timex.ToLayout("c"))
	assert.Eq(t, time.RFC3339, timex.ToLayout("Y-m-dTH:I:SP"))
}
