package nightstrike

const (
	backgroundColorHandle    = 0
	backgroundBitsHandle     = 1
	welcomeDisplayfontHandle = 2
	infofontHandle           = 11

	heliWidth        = 100
	heliHeight       = 62
	soldierRunWidth  = 34
	soldierRunHeight = 40
)

var (
	defensorFrontShape  = shape{2, 104, 64, 0}
	defensorTurretShape = shape{3, 38, 50, 50}
	missileAShape       = shape{4, 17, 51, 51}
	missileCShape       = shape{5, 19, 35, 35}
	heliShape           = shape{6, 100, 62, 0}
	copterFallShape     = shape{7, 45, 37, 0}
	fireShape           = shape{8, 26, 30, 0}
	explodeBigShape     = shape{9, 100, 70, 100}
	soldierRunShape     = shape{10, 34, 40, 0}
)

var assetsLevels = [...]string{
	assetsLevel0,
	assetsLevel1,
	assetsLevel2,
	//assetsLevel3,
	//assetsLevel4,
}
