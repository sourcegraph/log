package log

type Actor struct {
	ActorUID      string
	Ip            string
	XForwardedFor string
}

//TODO describe the 'actor takes an action on entity' idea

func (z *zapAdapter) Audit(actor Actor, action string, fields ...Field) {
	fields = append(fields, String("audit", "true"))

	fields = append(fields, Object("audit.actor",
		String("actorUID", actor.ActorUID),
		String("ip", actor.Ip),
		String("X-Forwarded-For", actor.XForwardedFor)))
	fields = append(fields, String("audit.action", action))
	fields = append(fields, String("audit.entity", z.fullScope))

	z.Info("audit action", fields...)
}
