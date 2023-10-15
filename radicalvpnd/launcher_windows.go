//go:build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func IsAdmin() bool {
	var sid *windows.SID

	//https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		log.Error(fmt.Sprintf("windows sid err: %s", err.Error()))
		return false
	}

	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		log.Error(fmt.Sprintf("token member error: %s", err.Error()))
		return false
	}

	return member
}
