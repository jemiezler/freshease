// Stub implementation for permission_handler on web
class Permission {
  static Permission activityRecognition = Permission();
  static Permission location = Permission();

  Future<PermissionStatus> request() async => PermissionStatus.denied;
}

enum PermissionStatus {
  granted,
  denied,
  restricted,
  limited,
  permanentlyDenied,
}
