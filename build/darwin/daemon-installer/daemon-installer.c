#include <syslog.h>
#include <string.h>
#include <stdio.h>
#include <ServiceManagement/ServiceManagement.h>
#include <Security/Security.h>

#define RADICAL_PACKAGE_NAME "com.radicalvpn.booter.helper"

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

    const char *prompt = "RadicalVPN requires a background service to establish VPN connections.\nIt is not possible to use RadicalVPN without the background service.\n\n";

    AuthorizationItem envItems = {kAuthorizationEnvironmentPrompt, strlen(prompt), (void *)prompt, 0};
    AuthorizationEnvironment env = { 1, &envItems };

    AuthorizationRef  authRef = NULL;
    OSStatus err = AuthorizationCreate(&authRights, &env, flags, &authRef);

    if(err == errAuthorizationSuccess)
    {
        //add the daemon starter the LaunchD
        bool success = SMJobBless(kSMDomainSystemLaunchd,
                      CFStringCreateWithCString(kCFAllocatorDefault, RADICAL_PACKAGE_NAME, kCFStringEncodingMacRoman),
                      (AuthorizationRef) authRef,
                      &error);

        AuthorizationFree(authRef, kAuthorizationFlagDefaults);

        if (success)
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

    return 1;
}

int get_installed_bundle_path(char* ret_buf, int buf_size) {
    CFStringRef helper_str = CFStringCreateWithCString(kCFAllocatorDefault, RADICAL_PACKAGE_NAME, kCFStringEncodingMacRoman);
    CFDictionaryRef retDict = SMJobCopyDictionary(kSMDomainSystemLaunchd, helper_str);
    if (helper_str != NULL)
    {
        CFRelease(helper_str);
    }

    if (retDict == NULL)
    {
        return 1;
    }

    int ret = 0;
    CFStringRef key = CFStringCreateWithCString(kCFAllocatorDefault, "ProgramArguments", kCFStringEncodingMacRoman);
    if (CFDictionaryContainsKey(retDict, key))
    {
        CFArrayRef program_arguments = CFDictionaryGetValue(retDict, key);
        if (program_arguments != NULL)
        {
            CFStringRef helperBundlePath = CFArrayGetValueAtIndex(program_arguments, 0);
            if (!CFStringGetCString(helperBundlePath, ret_buf, buf_size, kCFStringEncodingMacRoman))
            {
                ret = 2;
            }
        }
        else
        {
            ret = 3;
        }
    }
    else
    {
        ret = 4;
    }

    if (key != NULL)
    {
        CFRelease(key);
    }

    if (retDict != NULL)
    {
        CFRelease(retDict);
    }

    return ret;
}

bool daemon_install_required() {
    char installed_bundle[256] = {0};

    if (get_installed_bundle_path(installed_bundle, 256) == 0)
    {
        printf("%s\n", installed_bundle);
        return false;
    }
    else {
        return true;
    }
}


int main(int argc, char **argv) {
    if (argc >= 2)
    {
        if (strcmp(argv[1], "--install") == 0)
        {
            printf("Installing daemon..");
            return install_daemon();
        }

        if (strcmp(argv[1], "--install-required") == 0)
        {
            bool install_required = daemon_install_required();
            install_required ? printf("Daemon install required") : printf("Daemon install not required");
            return install_required ? 0 : 1;
        }
    }

    // no matching option
    printf("Invalid usage.\n");
    printf("Arguments:\n");
    printf("    --install\n");

    return 1;
}