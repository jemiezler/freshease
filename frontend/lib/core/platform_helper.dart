import 'package:flutter/foundation.dart';
import 'dart:io' if (dart.library.html) 'constants/platform.dart';

class PlatformHelper {
  static bool get isAndroid {
    // This check is now slightly redundant, since Platform.isAndroid
    // will be 'false' on web anyway, but it's good for clarity.
    if (kIsWeb) return false;

    // This now works!
    // On mobile, it calls dart:io.Platform.isAndroid
    // On web, it calls platform_helper_stub.Platform.isAndroid
    return Platform.isAndroid;
  }

  static bool get isIOS {
    if (kIsWeb) return false;
    return Platform.isIOS;
  }

  static bool get isMobile => isAndroid || isIOS;
}
