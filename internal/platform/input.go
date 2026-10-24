package platform

// Input is the browser control state sampled by one Go update.
type Input struct{ Left, Right, Up, Down, Attack, Use bool }
