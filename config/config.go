package config

import "time"

var SecretKey = []byte("topSecretlike")

var ExpJWT = time.Hour * 72
var SaltForRefreshToken = "asf21rwqASFSE@#!@WRFGE!@!RSFsafwqrq21rfwefsacasgw3E~!QEW"
var ExpireRefreshToken = time.Hour * 24 * 31
