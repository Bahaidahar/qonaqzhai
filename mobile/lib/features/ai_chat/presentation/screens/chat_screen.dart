import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../core/i18n/i18n.dart';
import '../../../../core/theme/app_theme.dart';
import '../viewmodels/chat_viewmodel.dart';
import '../widgets/ai_block.dart';

class ChatScreen extends ConsumerStatefulWidget {
  const ChatScreen({super.key});

  @override
  ConsumerState<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends ConsumerState<ChatScreen> {
  final _input = TextEditingController();
  final _scroll = ScrollController();

  @override
  void dispose() {
    _input.dispose();
    _scroll.dispose();
    super.dispose();
  }

  void _send() {
    final text = _input.text.trim();
    if (text.isEmpty) return;
    _input.clear();
    ref.read(chatProvider.notifier).send(text);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scroll.hasClients) {
        _scroll.animateTo(
          _scroll.position.maxScrollExtent,
          duration: const Duration(milliseconds: 200),
          curve: Curves.easeOut,
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(chatProvider);
    final p = AppPalette.of(context);
    final isEmpty = state.messages.isEmpty && !state.thinking;

    return Scaffold(
      appBar: AppBar(
        leading: Builder(
          builder: (ctx) => IconButton(
            icon: const Icon(CupertinoIcons.sidebar_left),
            tooltip: tr(ref, 'chat_history'),
            onPressed: () => Scaffold.of(ctx).openDrawer(),
          ),
        ),
        title: Text(tr(ref, 'nav_chat'),
            style: GoogleFonts.manrope(fontWeight: FontWeight.w700, fontSize: 17)),
        actions: [
          IconButton(
            icon: const Icon(CupertinoIcons.bell),
            tooltip: tr(ref, 'nav_notifications'),
            onPressed: () => context.push('/notifications'),
          ),
          IconButton(
            icon: const Icon(CupertinoIcons.add),
            tooltip: tr(ref, 'chat_new'),
            onPressed: () => ref.read(chatProvider.notifier).newChat(),
          ),
        ],
      ),
      drawer: const _ChatHistoryDrawer(),
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: isEmpty
                  ? _ChatHero(onPickPreset: (text) {
                      _input.text = text;
                      _send();
                    })
                  : ListView.builder(
                      controller: _scroll,
                      padding: const EdgeInsets.fromLTRB(16, 16, 16, 24),
                      itemCount: state.messages.length + (state.thinking ? 1 : 0),
                      itemBuilder: (_, i) {
                        if (i == state.messages.length) {
                          return Padding(
                            padding: const EdgeInsets.symmetric(vertical: 8),
                            child: _ThinkingDots(color: p.mutedFg),
                          );
                        }
                        final m = state.messages[i];
                        final isUser = m.role == 'user';
                        return Align(
                          alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
                          child: Padding(
                            padding: const EdgeInsets.symmetric(vertical: 6),
                            child: Column(
                              crossAxisAlignment:
                                  isUser ? CrossAxisAlignment.end : CrossAxisAlignment.start,
                              children: [
                                ConstrainedBox(
                                  constraints: BoxConstraints(
                                    maxWidth: MediaQuery.of(context).size.width * 0.82,
                                  ),
                                  child: Container(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 14, vertical: 10),
                                    decoration: BoxDecoration(
                                      color: isUser ? p.primary : p.card,
                                      borderRadius: BorderRadius.circular(16),
                                      border: isUser ? null : Border.all(color: p.border),
                                    ),
                                    child: Text(
                                      m.text,
                                      style: GoogleFonts.manrope(
                                        fontSize: 14,
                                        height: 1.45,
                                        color: isUser ? p.onPrimary : p.fg,
                                      ),
                                    ),
                                  ),
                                ),
                                for (final b in m.blocks)
                                  Padding(
                                    padding: const EdgeInsets.only(top: 8),
                                    child: AiBlock(block: b),
                                  ),
                              ],
                            ),
                          ),
                        );
                      },
                    ),
            ),
            if (state.error != null)
              Padding(
                padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                  decoration: BoxDecoration(
                    color: p.destructive.withValues(alpha: 0.08),
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(color: p.destructive.withValues(alpha: 0.3)),
                  ),
                  child: Text(
                    state.error!,
                    style: GoogleFonts.manrope(
                      fontSize: 12,
                      color: p.destructive,
                    ),
                  ),
                ),
              ),
            _Composer(
              controller: _input,
              thinking: state.thinking,
              onSend: _send,
            ),
          ],
        ),
      ),
    );
  }
}

class _Composer extends StatelessWidget {
  const _Composer({
    required this.controller,
    required this.thinking,
    required this.onSend,
  });

  final TextEditingController controller;
  final bool thinking;
  final VoidCallback onSend;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Padding(
      padding: const EdgeInsets.fromLTRB(12, 4, 12, 12),
      child: Container(
        padding: const EdgeInsets.fromLTRB(6, 6, 6, 6),
        decoration: BoxDecoration(
          color: p.card,
          borderRadius: BorderRadius.circular(18),
          border: Border.all(color: p.border),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.end,
          children: [
            Expanded(
              child: TextField(
                controller: controller,
                minLines: 1,
                maxLines: 5,
                style: GoogleFonts.manrope(fontSize: 14, color: p.fg),
                decoration: InputDecoration(
                  filled: false,
                  hintText: 'Опиши событие, бюджет, локацию…',
                  hintStyle: GoogleFonts.manrope(fontSize: 14, color: p.mutedFg),
                  border: InputBorder.none,
                  enabledBorder: InputBorder.none,
                  focusedBorder: InputBorder.none,
                  contentPadding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
                ),
                onSubmitted: (_) => onSend(),
              ),
            ),
            SizedBox(
              width: 40,
              height: 40,
              child: Material(
                color: thinking ? p.muted : p.primary,
                shape: const CircleBorder(),
                child: InkWell(
                  customBorder: const CircleBorder(),
                  onTap: thinking ? null : onSend,
                  child: Icon(
                    CupertinoIcons.arrow_up,
                    size: 18,
                    color: thinking ? p.mutedFg : p.onPrimary,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ChatHero extends ConsumerWidget {
  const _ChatHero({required this.onPickPreset});
  final void Function(String) onPickPreset;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final p = AppPalette.of(context);
    final presets = <String>[
      'Свадьба на 120 гостей, бюджет 5 млн ₸, Алматы — собери план',
      'Корпоратив на 80 человек, нужен фотограф и кейтеринг',
      'День рождения ребёнка 6 лет, площадка с аниматорами',
    ];
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: p.primary.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(14),
              border: Border.all(color: p.primary.withValues(alpha: 0.25)),
            ),
            child: Icon(CupertinoIcons.sparkles, size: 22, color: p.primary),
          ),
          const SizedBox(height: 16),
          Text(
            'Планируй событие\nкак сообщение',
            style: GoogleFonts.manrope(
              fontSize: 36,
              fontWeight: FontWeight.w700,
              color: p.fg,
              letterSpacing: -1.4,
              height: 1.05,
            ),
          ),
          const SizedBox(height: 10),
          Text(
            'Опиши идею — qonaqzhai подберёт исполнителей, тайминг и бюджет.',
            style: GoogleFonts.manrope(fontSize: 14, color: p.mutedFg, height: 1.5),
          ),
          const SizedBox(height: 28),
          for (final preset in presets) ...[
            _PresetCard(text: preset, onTap: () => onPickPreset(preset)),
            const SizedBox(height: 10),
          ],
        ],
      ),
    );
  }
}

class _PresetCard extends StatelessWidget {
  const _PresetCard({required this.text, required this.onTap});
  final String text;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            color: p.card,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(color: p.border),
          ),
          child: Row(
            children: [
              Expanded(
                child: Text(
                  text,
                  style: GoogleFonts.manrope(fontSize: 13.5, color: p.fg, height: 1.45),
                ),
              ),
              Icon(CupertinoIcons.arrow_up_right, size: 16, color: p.mutedFg),
            ],
          ),
        ),
      ),
    );
  }
}

class _ThinkingDots extends StatefulWidget {
  const _ThinkingDots({required this.color});
  final Color color;

  @override
  State<_ThinkingDots> createState() => _ThinkingDotsState();
}

class _ThinkingDotsState extends State<_ThinkingDots>
    with SingleTickerProviderStateMixin {
  late final AnimationController _ctl =
      AnimationController(vsync: this, duration: const Duration(milliseconds: 1200))
        ..repeat();

  @override
  void dispose() {
    _ctl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 18,
      child: AnimatedBuilder(
        animation: _ctl,
        builder: (_, __) {
          return Row(
            mainAxisSize: MainAxisSize.min,
            children: List.generate(3, (i) {
              final t = ((_ctl.value + i * 0.2) % 1.0);
              final size = 4 + 4 * (1 - (2 * t - 1).abs());
              return Padding(
                padding: const EdgeInsets.symmetric(horizontal: 3),
                child: Container(
                  width: size,
                  height: size,
                  decoration: BoxDecoration(
                    color: widget.color,
                    shape: BoxShape.circle,
                  ),
                ),
              );
            }),
          );
        },
      ),
    );
  }
}

class _ChatHistoryDrawer extends ConsumerStatefulWidget {
  const _ChatHistoryDrawer();

  @override
  ConsumerState<_ChatHistoryDrawer> createState() => _ChatHistoryDrawerState();
}

class _ChatHistoryDrawerState extends ConsumerState<_ChatHistoryDrawer> {
  List<Map<String, dynamic>> _items = const [];
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      _items = await ref.read(chatRepositoryProvider).listChats();
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final p = AppPalette.of(context);
    return Drawer(
      backgroundColor: p.bg,
      child: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 18, 12, 12),
              child: Row(
                children: [
                  Expanded(
                    child: Text(
                      tr(ref, 'chat_history'),
                      style: GoogleFonts.manrope(
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                        color: p.fg,
                        letterSpacing: -0.3,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(CupertinoIcons.add),
                    onPressed: () {
                      Navigator.pop(context);
                      ref.read(chatProvider.notifier).newChat();
                    },
                  ),
                ],
              ),
            ),
            Divider(height: 1, color: p.border),
            Expanded(
              child: _loading
                  ? Center(
                      child: Text(
                        tr(ref, 'common_loading'),
                        style: GoogleFonts.manrope(color: p.mutedFg),
                      ),
                    )
                  : _items.isEmpty
                      ? Center(
                          child: Text(
                            tr(ref, 'common_empty'),
                            style: GoogleFonts.manrope(color: p.mutedFg),
                          ),
                        )
                      : ListView.separated(
                          padding: const EdgeInsets.symmetric(vertical: 8),
                          itemCount: _items.length,
                          separatorBuilder: (_, __) => const SizedBox(height: 2),
                          itemBuilder: (_, i) {
                            final c = _items[i];
                            return ListTile(
                              tileColor: Colors.transparent,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(10),
                              ),
                              title: Text(
                                c['title']?.toString() ?? '…',
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                                style: GoogleFonts.manrope(
                                  fontSize: 13.5,
                                  color: p.fg,
                                ),
                              ),
                              trailing: IconButton(
                                icon: Icon(CupertinoIcons.trash, size: 18, color: p.mutedFg),
                                onPressed: () async {
                                  await ref
                                      .read(chatRepositoryProvider)
                                      .deleteChat(c['id'].toString());
                                  await _load();
                                },
                              ),
                              onTap: () {
                                Navigator.pop(context);
                                ref
                                    .read(chatProvider.notifier)
                                    .loadChat(c['id'].toString());
                              },
                            );
                          },
                        ),
            ),
          ],
        ),
      ),
    );
  }
}
