package main

var _clickEffect IAnimate

func GetClickEffect() IAnimate {
	if _clickEffect == nil {
		_clickEffect = NewAnimateOneImage(GetMainView().window.Render, "assets/cursor.png", 8, 100, 11, 30)
	}
	return _clickEffect
}
