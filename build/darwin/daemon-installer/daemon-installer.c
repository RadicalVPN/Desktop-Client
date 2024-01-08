#include <Security/Security.h>
#include <ServiceManagement/ServiceManagement.h>
#include <stdio.h>
#include <string.h>
#include <syslog.h>

#define RADICAL_PACKAGE_NAME "com.radicalvpn.booter.helper"
#define INSTALLED_HELPER_BINARY "/Library/PrivilegedHelperTools/com.radicalvpn.booter.helper"
#define BUNDLE_HELPER_BINARY "/Applications/RadicalVPN.app/Contents/MacOS/RadicalVPN-Installer.app/Contents/Library/LaunchServices/com.radicalvpn.booter.helper"
#define LAUNCH_DAEMON_PLIST "/Library/LaunchDaemons/com.radicalvpn.booter.helper.plist"

void log_error_ref(CFErrorRef err)
{
	if (err == NULL)
		return;

	CFStringRef error_txt = CFErrorCopyDescription(err);
	if (error_txt == NULL)
		return;

	const char *ptr = CFStringGetCStringPtr(error_txt, kCFStringEncodingUTF8);
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

CFDictionaryRef get_bundle_dictionary(const char *bundle_path)
{
	CFStringRef bundle_str = CFStringCreateWithCString(kCFAllocatorDefault, bundle_path, kCFStringEncodingMacRoman);
	CFStringRef bundle_str_escaped = CFURLCreateStringByAddingPercentEscapes(NULL, bundle_str, NULL, NULL, kCFStringEncodingUTF8);
	CFURLRef url = CFURLCreateWithString(NULL, bundle_str_escaped, NULL);
	CFDictionaryRef dic = CFBundleCopyInfoDictionaryForURL(url);

	if (url != NULL)
	{
		CFRelease(url);
	}

	if (bundle_str != NULL)
	{
		CFRelease(bundle_str);
	}

	return dic;
}

int get_bundle_version(const char *bundle_path, char *ret_buf, int buf_size)
{
	CFDictionaryRef dict = get_bundle_dictionary(bundle_path);
	if (dict == NULL)
	{
		return 1;
	}

	int ret = 0;
	CFStringRef key = CFStringCreateWithCString(kCFAllocatorDefault, "CFBundleVersion", kCFStringEncodingMacRoman);

	if (dict != NULL && CFDictionaryContainsKey(dict, key))
	{
		CFStringRef ver = CFDictionaryGetValue(dict, key);

		if (ver == NULL)
		{
			ret = 2;
		}
		else
		{
			if (!CFStringGetCString(ver, ret_buf, buf_size, kCFStringEncodingMacRoman))
			{
				ret = 3;
			}

			CFRelease(ver);
		}
	}
	else
		ret = 4;

	if (key != NULL)
	{
		CFRelease(key);
	}

	return ret;
}

int remove_daemon(AuthorizationRef authRef)
{
	CFErrorRef error = NULL;

	bool success = SMJobRemove(kSMDomainSystemLaunchd, CFStringCreateWithCString(kCFAllocatorDefault, RADICAL_PACKAGE_NAME, kCFStringEncodingMacRoman), authRef, true, &error);
	if (!success)
	{
		if (error != NULL)
		{
			log_error_ref(error);
			CFRelease(error);
		}
		return 1;
	}

	char *files[] = {LAUNCH_DAEMON_PLIST, INSTALLED_HELPER_BINARY, NULL};
	OSStatus err = AuthorizationExecuteWithPrivileges(authRef, (const char *)"/bin/rm", kAuthorizationFlagDefaults, files, NULL);
	if (err)
	{
		return 2;
	}

	return 0;
}

int install_daemon()
{
	bool upgrade = false;
	char current_version[64] = {0};
	char installed_version[64] = {0};
	int ret1 = get_bundle_version(BUNDLE_HELPER_BINARY, current_version, sizeof(current_version));
	int ret2 = get_bundle_version(INSTALLED_HELPER_BINARY, installed_version, sizeof(installed_version));

	if (installed_version[0] != 0)
	{
		if (strcmp(installed_version, current_version) == 0)
		{
			return 1;
		}

		upgrade = true;
	}

	CFErrorRef error = NULL;
	AuthorizationItem authItem = {kSMRightBlessPrivilegedHelper, 0, NULL, 0};
	AuthorizationRights authRights = {1, &authItem};
	AuthorizationFlags flags = kAuthorizationFlagDefaults | kAuthorizationFlagInteractionAllowed | kAuthorizationFlagPreAuthorize | kAuthorizationFlagExtendRights;

	const char *prompt_upgrading = "A new version of the RadicalVPN daemon has been installed.";
	const char *prompt = "RadicalVPN requires a background service to establish VPN connections.\nIt is not possible to use RadicalVPN without the background service.\n\n";

	if (upgrade)
	{
		prompt = prompt_upgrading;
	}

	AuthorizationItem envItems = {kAuthorizationEnvironmentPrompt, strlen(prompt), (void *)prompt, 0};
	AuthorizationEnvironment env = {1, &envItems};

	AuthorizationRef authRef = NULL;
	OSStatus err = AuthorizationCreate(&authRights, &env, flags, &authRef);

	if (err == errAuthorizationSuccess)
	{
		if (upgrade)
		{
			// remove the old daemon
			remove_daemon(authRef);
		}

		// add the daemon starter the LaunchD
		bool success = SMJobBless(kSMDomainSystemLaunchd, CFStringCreateWithCString(kCFAllocatorDefault, RADICAL_PACKAGE_NAME, kCFStringEncodingMacRoman), (AuthorizationRef)authRef, &error);

		AuthorizationFree(authRef, kAuthorizationFlagDefaults);

		if (success)
		{
			printf("Daemon installed.");
			return 0;
		}
		else
		{
			log_error_ref(error);

			if (error != NULL)
			{
				CFRelease(error);
			}
			return 1;
		}
	}

	return 1;
}

bool daemon_install_required()
{
	char current_version[64] = {0};
	char installed_version[64] = {0};

	int ret1 = get_bundle_version(BUNDLE_HELPER_BINARY, current_version, sizeof(current_version));
	int ret2 = get_bundle_version(INSTALLED_HELPER_BINARY, installed_version, sizeof(installed_version));

	if (ret1 == 0 && ret2 == 0)
	{
		if (strcmp(current_version, installed_version) == 0)
		{
			return false;
		}
		else
		{
			return true;
		}
	}
	else
	{
		return false;
	}
}

int main(int argc, char **argv)
{
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
	printf("    --install-required\n");

	return 1;
}