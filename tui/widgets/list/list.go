package list

// Until termdash allows for non-internal widgets, this is a no-go.

type List struct {
}

// func New(opts ...Option) (*List, error) {
// 	return nil, nil
// }

// // When the infrastructure calls Draw(), the widget must block on the call
// // until it finishes drawing onto the provided canvas. When given the
// // canvas, the widget must first determine its size by calling
// // Canvas.Size(), then limit all its drawing to this area.
// //
// // The widget must not assume that the size of the canvas or its content
// // remains the same between calls.
// func (l *List) Draw(cvs *canvas.Canvas) error {
// 	return nil
// }

// // Keyboard is called when the widget is focused on the dashboard and a key
// // shortcut the widget registered for was pressed. Only called if the widget
// // registered for keyboard events.
// func (l *List) Keyboard(k *terminalapi.Keyboard) error {
// 	return nil
// }

// // Mouse is called when the widget is focused on the dashboard and a mouse
// // event happens on its canvas. Only called if the widget registered for mouse
// // events.
// func (l *List) Mouse(m *terminalapi.Mouse) error {
// 	return nil
// }

// // Options returns registration options for the widget.
// // This is how the widget indicates to the infrastructure whether it is
// // interested in keyboard or mouse shortcuts, what is its minimum canvas
// // size, etc.
// //
// // Most widgets will return statically compiled options (minimum and
// // maximum size, etc.). If the returned options depend on the runtime state
// // of the widget (e.g. the user data provided to the widget), the widget
// // must protect against a case where the infrastructure calls the Draw
// // method with a canvas that doesn't meet the requested options. This is
// // because the data in the widget might change between calls to Options and
// // Draw.
// func (l *List) Options() Options {
// 	return nil
// }
