package main

import (
	. "github.com/SilentGopherLnx/easygolang"
	. "github.com/SilentGopherLnx/easygolang/easygtk"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var select_mode bool
var select_x1, select_y1, select_x2, select_y2 int

func FileSelector_GetList() []string {
	arr := []string{}
	for j := 0; j < len(arr_blocks); j++ {
		if arr_blocks[j].check.GetActive() {
			arr = append(arr, arr_blocks[j].fname) // FolderPathEndSlash(arr_blocks[j].fpath)+
		}
	}
	return arr
}

func FileSelector_Draw(dy int, ctx *cairo.Context) {
	if select_x1 > 0 && select_y1 > 0 && select_x2 > 0 && select_y2 > 0 {
		c := GTK_ColorOfSelected()
		ctx.SetSourceRGBA(c[0], c[1], c[2], 1.0) //0.4, 0.7, 0.8, 1.0) // BLUE DARK
		ctx.Rectangle(float64(select_x1), float64(select_y1-dy), float64(select_x2-select_x1), float64(select_y2-select_y1))
		ctx.Fill()
	}
}

func FileSelector_MouseAtSelectZone(x0, y0 int) bool {
	at_zone := false
	for j := 0; j < len(arr_blocks); j++ {
		at_zone = at_zone || arr_blocks[j].IsClickedIn(x0, y0)
	}
	return !at_zone
}

func FileSelector_MousePressed(event *gdk.Event, scroll *gtk.ScrolledWindow) (int, int, int, bool) {
	mousekey, x1, y1 := GTK_MouseKeyOfEvent(event)
	_, dy := GTK_ScrollGetValues(scroll)
	y1 += dy
	zone := FileSelector_MouseAtSelectZone(x1, y1)
	if mousekey == 1 && zone {
		select_x1 = x1
		select_y1 = y1
		select_x2 = 0
		select_y2 = 0
		Prln("select mouse1_down " + I2S(x1) + "/" + I2S(y1+dy))
	}
	return mousekey, x1, y1, zone
}

func FileSelector_MouseMoved(event *gdk.Event, scroll *gtk.ScrolledWindow, redraw func()) {
	if select_x1 > 0 && select_y1 > 0 {
		_, x2, y2 := GTK_MouseKeyOfEvent(event)
		_, dy := GTK_ScrollGetValues(scroll)
		y2 += dy
		select_x2 = x2
		select_y2 = y2
		//Prln("rect " + I2S(select_x1) + "," + I2S(select_y1) + " / " + I2S(select_x2) + "," + I2S(select_y2))
		for j := 0; j < len(arr_blocks); j++ {
			is_inside := arr_blocks[j].IsInSelectRect(select_x1, select_y1, select_x2, select_y2)
			arr_blocks[j].SetSelected(is_inside)
		}
		redraw()
	}
}

func FileSelector_MouseRelease(event *gdk.Event, scroll *gtk.ScrolledWindow, redraw func()) {
	mousekey, x2, y2 := GTK_MouseKeyOfEvent(event)
	_, dy := GTK_ScrollGetValues(scroll)
	y2 += dy
	zone2 := FileSelector_MouseAtSelectZone(x2, y2)
	if mousekey == 1 {
		if select_x1 == x2 && select_y1 == y2 && zone2 {
			Prln("select mouse1_up with reset")
			FileSelector_ResetChecks()
		} else {
			Prln("select mouse1_up")
		}
		FileSelector_ResetRect()
		redraw()
	}
}

func FileSelector_ResetRect() {
	select_x1 = 0
	select_y1 = 0
	select_x2 = 0
	select_y2 = 0
}

func FileSelector_ResetChecks() {
	for j := 0; j < len(arr_blocks); j++ {
		arr_blocks[j].SetSelected(false)
	}
}
