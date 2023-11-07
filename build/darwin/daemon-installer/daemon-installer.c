#include <syslog.h>
#include <string.h>
#include <stdio.h>
#include <ServiceManagement/ServiceManagement.h>
#include <Security/Security.h>

#define RADICAL_PACKAGE_NAME "com.radicalvpn.booter"

void log_error_ref(CFErrorRef err) {
	if (err == NULL)
		return;

    CFStringRef error_txt = CFErrorCopyDescription(err);
    if (error_txt == NULL)
        return;

	const char  *ptr      = CFStringGetCStringPtr(error_txt, kCFStringEncodingUTF8);
    if (ptr == NULL)
    {
        const CFIndex bufSize = 1024;
        char buf[bufSize];

        if (CFStringGetCString(error_txt, buf, bufSize, kCFStringEncodingUTF8))
        {
            printf("%s", buf);
        }
        return;
    }

    printf("%s", ptr);
}

int install_daemon() {

    CFErrorRef error = NULL;
    AuthorizationItem authItem = { kSMRightBlessPrivilegedHelper, 0, NULL, 0 };
    AuthorizationRights authRights = { 1, &authItem };
    AuthorizationFlags flags = kAuthorizationFlagDefaults |
                               kAuthorizationFlagInteractionAllowed |
                               kAuthorizationFlagPreAuthorize |
                               kAuthorizationFlagExtendRights;

    const char *prompt = "Penis.\n\n";


    AuthorizationItem envItems = {kAuthorizationEnvironmentPrompt, strlen(prompt), (void *)prompt, 0};
    AuthorizationEnvironment env = { 1, &envItems };

    AuthorizationRef  authRef = NULL;
    OSStatus err = AuthorizationCreate(&authRights, &env, flags, &authRef);

    if(err == errAuthorizationSuccess)
    {
        //add the daemon starter the LaunchD
        bool isSuccess = SMJobBless(kSMDomainSystemLaunchd,
                      CFStringCreateWithCString(kCFAllocatorDefault, RADICAL_PACKAGE_NAME, kCFStringEncodingMacRoman),
                      (AuthorizationRef) authRef,
                      &error);

        AuthorizationFree(authRef, kAuthorizationFlagDefaults);

        if (isSuccess)
        {
            printf("Daemon installed.");
            return 0;
        }
        else
        {
            log_error_ref(error);
            printf("%s", &error);

            if (error != NULL) {
                CFRelease(error);
            }
            return 1;
        }
    }


    printf("err: %d\n", err);
}

int main(int argc, char **argv) {
    install_daemon();
}