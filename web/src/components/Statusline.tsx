export default function Statusline() {
  return (
    <footer className="flex flex-none items-center justify-between gap-4 border-t border-frame-dim bg-black px-3.5 py-1 text-[10.5px] tracking-[.06em] text-dimmer">
      <span className="whitespace-nowrap [&_kbd]:mr-1 [&_kbd]:bg-rule [&_kbd]:px-1.5 [&_kbd]:font-[inherit] [&_kbd]:text-fg">
        <kbd>1</kbd>
        <kbd>2</kbd>
        <kbd>3</kbd>
        switch view
      </span>
      <span className="whitespace-nowrap text-frame">— NORMAL —</span>
    </footer>
  );
}
