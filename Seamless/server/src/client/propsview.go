package main

import (
	"fmt"
)

const (
	PVXBegin = 30
	PVYBegin = 55

	GridWidth  = 160
	GridHeight = 25
)

type PropsView struct {
	mv *MainView
}

func (pv *PropsView) DrawProps() {
	props := &GetClient().UserView.props
	pv.drawOne(0, 0, fmt.Sprintf("%10v%-10v", "UserName: ", GetClient().GetUserName()))
	pv.drawOne(1, 0, fmt.Sprintf("%10v%-10v", "EntityID: ", GetClient().User.EntityID))
	pv.drawOne(2, 0, fmt.Sprintf("%10v%-10v", "UID: ", GetClient().UID))
	pv.drawOne(3, 0, fmt.Sprintf("%10v%-10v", "Hp: ", props.Hp))
	pv.drawOne(4, 0, fmt.Sprintf("%10v%-10v", "MaxHp: ", props.MaxHp))
	pv.drawOne(5, 0, fmt.Sprintf("%10v%-10v", "Attack: ", props.Attack))
	pv.drawOne(6, 0, fmt.Sprintf("%10v%-10v", "Defence: ", props.Defence))

	us := &GetClient().UserView.us
	pv.drawOne(8, 0, fmt.Sprintf("%10v%-10.2f", "X: ", us.pos.X))
	pv.drawOne(9, 0, fmt.Sprintf("%10v%-10.2f", "Z: ", us.pos.Z))

	cellSta := "ok"
	if !GetClient().IsCellOk() {
		cellSta = "lost connect"
	}
	pv.drawOne(0, 5, "CellState:"+cellSta)
}

func (pv *PropsView) drawOne(i, j int32, text string) {
	pv.mv.DrawText(PVXBegin+j*GridWidth, PVYBegin+i*GridHeight, text, 0, 245, 255)
}
