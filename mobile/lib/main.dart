import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'core/i18n/i18n.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'core/theme/theme_controller.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  // Force the semantic tree on so Maestro / VoiceOver can see widget labels.
  // No-op for production users; cheap on the rendering pipeline.
  RendererBinding.instance.ensureSemantics();
  runApp(const ProviderScope(child: QonaqzhaiApp()));
}

class QonaqzhaiApp extends ConsumerWidget {
  const QonaqzhaiApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);
    final locale = ref.watch(localeProvider);
    final themePref = ref.watch(themeControllerProvider);
    return MaterialApp.router(
      title: 'Qonaqzhai',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light(),
      darkTheme: AppTheme.dark(),
      themeMode: themePref.materialMode,
      routerConfig: router,
      locale: locale.flutterLocale,
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [Locale('kk'), Locale('ru'), Locale('en')],
    );
  }
}
