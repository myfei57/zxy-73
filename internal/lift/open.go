package lift

type SteamVentPort interface {
	Vent() error
}

func (l *LiftControl) Open(steam SteamVentPort) error {
	// 泄压必须先于开模执行：否则模具在带压状态下打开，
	// 高温蒸汽会外泄，存在安全风险。
	if err := steam.Vent(); err != nil {
		return err
	}
	l.closed = false
	return nil
}
