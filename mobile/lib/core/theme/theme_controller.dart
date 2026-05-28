import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

enum AppThemePref { light, dark, system }

extension AppThemePrefX on AppThemePref {
  String get wire => switch (this) {
        AppThemePref.light => 'light',
        AppThemePref.dark => 'dark',
        AppThemePref.system => 'system',
      };

  ThemeMode get materialMode => switch (this) {
        AppThemePref.light => ThemeMode.light,
        AppThemePref.dark => ThemeMode.dark,
        AppThemePref.system => ThemeMode.system,
      };

  static AppThemePref parse(String? s) {
    switch (s) {
      case 'light':
        return AppThemePref.light;
      case 'dark':
        return AppThemePref.dark;
      default:
        return AppThemePref.system;
    }
  }
}

/// Persists the user's appearance preference to SharedPreferences using the
/// same `qonaqzhai_theme` key the web app reads — so a future shared profile
/// sync would just work.
class ThemeController extends StateNotifier<AppThemePref> {
  ThemeController() : super(AppThemePref.system) {
    _restore();
  }

  static const _key = 'qonaqzhai_theme';

  Future<void> _restore() async {
    final p = await SharedPreferences.getInstance();
    state = AppThemePrefX.parse(p.getString(_key));
  }

  Future<void> set(AppThemePref pref) async {
    state = pref;
    final p = await SharedPreferences.getInstance();
    await p.setString(_key, pref.wire);
  }
}

final themeControllerProvider =
    StateNotifierProvider<ThemeController, AppThemePref>((ref) {
  return ThemeController();
});
