package bit

import (
	"bufio"
	"errors"
	"os"
	"fmt"
	"unicode"
	"github.com/jroimartin/gocui"
)




const BitUp = "▲"
const BitDown = "▼"
const BitLeft = "◀"
const BitRight = "▶"

const HelpText = "n: next step    p: previous step    f: first step    l: last step    s: switch world    q: quit"




/*
   Enumeration to represent the various colors of a square
*/
type Color int
const (
	White Color = iota
	Black
	Red
	Blue
	Green
)
/*
   Deletes the view that shows the steps to replace it with a new one that has been made wide enough
   to fit the width of the current step.
*/
func update_step_view() {
	gui_i.DeleteView("Steps")
	if step_v, err := gui_i.SetView("Steps", max_x/2-len(bit_state_actions[current_state]),max_y/2-12,max_x/2+len(bit_state_actions[current_state]), max_y/2-10); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(step_v, current_state, ": ", bit_state_actions[current_state])
	}
}
/*
   Goes to the initial step of the bit world
*/
func first_state(gui *gocui.Gui, world *gocui.View) error {
	world.Clear()
	current_state = 0
	print_world(world, bit_states[current_state].face, bit_states[current_state].world)
	update_step_view()
	return nil
}

/*
   Goes to the next step of the bit world.
*/
func next_state(gui *gocui.Gui, world *gocui.View) error {
	if current_state < len(bit_states) - 1 {
		world.Clear()
		current_state++
		print_world(world, bit_states[current_state].face, bit_states[current_state].world)
		update_step_view()
		return nil
	}
	return nil
}

/*
   Goes to the last step of the bit world.
*/
func last_state(gui *gocui.Gui, world *gocui.View) error {
	world.Clear()
	current_state = len(bit_states) - 1
	print_world(world, bit_states[current_state].face, bit_states[current_state].world)
	update_step_view()
	return nil
}

/*
   Goes to the previous step of the bit world.
*/
func prev_state(gui *gocui.Gui, world *gocui.View) error {
	
	if current_state > 0 {
		world.Clear()
		current_state--
		print_world(world, bit_states[current_state].face, bit_states[current_state].world)
		update_step_view()
		return nil
	}
	return nil
}

/*
   Determines if the we are showing what the final world should look like or what the current step of bit looks like.
*/
var final_world bool = false

/*
   Switches between showing the answer and the current step of bit.
*/
func switch_world(gui *gocui.Gui, world *gocui.View) error {
	world.Clear()
	final_world = !final_world
	if final_world {
		print_world(world, bit_states[current_state].face, bit_states[current_state].final_world)
	} else {
		print_world(world, bit_states[current_state].face, bit_states[current_state].world)
	}
	return nil
}

/*
   Function that quits the bit gui.
*/
func quit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

/*
   This function resets the global variable that are used to keep track of the gui.
*/
func reset_gui_globals() {
	gui_i = nil
	max_x = 0
	max_y = 0
}
// This is a pointer to the gui that is being used to display the bit world.
var gui_i *gocui.Gui = nil
//This represents the width of the gui.
var max_x int = 0
//This represents the height of the gui.
var max_y int = 0

/*
   This function draws the bit world in the gui and sets the keybindings for it.
*/
func setup_world_view(gui *gocui.Gui, x0 int, y0 int, x1 int, y1 int) error {
	
	if world, err := gui.SetView("World", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		current_state = len(bit_states) - 1

		print_world(world, bit_states[current_state].face, bit_states[current_state].world)

		if err := gui.SetKeybinding("World", rune('n'), gocui.ModNone, next_state); err != nil {
			return err
		}
		if err := gui.SetKeybinding("World", rune('p'), gocui.ModNone, prev_state); err != nil {
			return err
		}
		if err := gui.SetKeybinding("World", rune('f'), gocui.ModNone, first_state); err != nil {
			return err
		}
		if err := gui.SetKeybinding("World", rune('l'), gocui.ModNone, last_state); err != nil {
			return err
		}
		if err := gui.SetKeybinding("", rune('q'), gocui.ModNone, quit); err != nil {
			return err
		}
		if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			return err
		}
		if err := gui.SetKeybinding("World", rune('s'), gocui.ModNone, switch_world); err != nil {
			return err
		}
		if _, err := gui.SetCurrentView("World"); err != nil {
			return err
		} 
	}
	return nil
}

/*
   This function initializes the step view for bit.
*/
func setup_step_view(gui *gocui.Gui, x0 int, y0 int, x1 int, y1 int) error {
	if step_view, err := gui.SetView("Steps", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(step_view, current_state, ": ", bit_state_actions[current_state])
	}
	return nil
}

/*
   This function initializes the help view for the bit gui.
*/
func setup_help_view(gui *gocui.Gui, x0 int, y0 int, x1 int, y1 int) error {
	
	if help_view, err := gui.SetView("Help", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(help_view, HelpText)
	}
	return nil
}

/*
   This function launches the gui for bit.
   
*/
func RunGui() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer gui.Close()
	defer reset_gui_globals()
	
	gui_i = gui

	max_x, max_y = gui.Size()

	var x0, y0, x1, y1 int
	
	if (world_width * 5) % 2 != 0 {
		x0 = max_x/2 - (world_width*5/2)
		x1 = max_x/2 + (world_width*5/2) +2
	} else {
		x0 = max_x/2 - (world_width*5/2)
		x1 = max_x/2 + (world_width*5/2) +1
	}
	if (world_height*3) % 2 != 0 {
		y0 = max_y/2 - (world_height*3/2)
		y1 = max_y/2 + (world_height*3/2) -1
	} else {
		y0 = max_y/2 - (world_height*3/2)
		y1 = max_y/2 + (world_height*3/2) -2
	}
		
	
	if err := setup_world_view(gui, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}


	if err := setup_step_view(gui, max_x/2-len(bit_state_actions[current_state]),max_y/2-12,max_x/2+len(bit_state_actions[current_state]), max_y/2-10); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if err := setup_help_view(gui, max_x/2 - len(HelpText)/2 -1, max_y/2+(world_height)+2, max_x/2+len(HelpText)/2 +1, max_y/2+(world_height) +4); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		os.Exit(1)
	}

}

/*
   This function resets the global variables that bit and the gui use to keep track of the state.
*/
func reset_bit_globals() {
	bit_snapshots = make([]int, 0)
	current_snapshot = 0
	current_state = 0
	bit_states = make([]*Bit, 1)
	bit_state_actions = make([]string, 1)
	error_occured = nil
	world_width = 0
	world_height = 0
}

// Snapshots of the world at steps the user has created.
var bit_snapshots []int = make([]int, 0)
var bit_snapshot_names []string = make([]string, 0)

// The current snapshot of the world.
var current_snapshot int = 0

// The current state of the world, used by the gui.
var current_state int = 0
// The states of the world at each step that bit has taken.
var bit_states []*Bit = make([]*Bit, 1)
// A list that is the action that bit made at a step
var bit_state_actions []string = make([]string, 1)

// A global that specifies that an error has occured so that we can skip every other bit step.
var error_occured error = nil

// The width of the world.
var world_width int = 0
// The height of the world.
var world_height int = 0

/*
   This function initializes the bit globals.
*/
func init_bit_state(bit *Bit) {
	dest := make([][]Square, len(bit.world))
	for i := range dest {
		dest[i] = make([]Square, len(bit.world[i]))
		copy(dest[i], bit.world[i])
	}
	bit_states[0] = &Bit{ face: bit.face, steps: bit.steps, x: bit.x, y: bit.y, world: dest, final_world: bit.final_world}
	bit_state_actions[0] = "initial state"
}

/*
   This function is called whenever bit does an action.
   This adds the action to the list of actions that bit has taken.
*/
func add_bit_state(bit *Bit, action string) {
	dest := make([][]Square, len(bit.world))
	for i := range dest {
		dest[i] = make([]Square, len(bit.world[i]))
		copy(dest[i], bit.world[i])
	}
	bit_states = append(bit_states, &Bit{ face: bit.face, steps: bit.steps, x: bit.x, y: bit.y, world: dest, final_world: bit.final_world})
	bit_state_actions = append(bit_state_actions, action)
}

/*
   This function is called whenever bit makes an invalid move.
   This function adds the error to the list of actions that bit has taken.
*/
func stop_display_error(msg string, bit *Bit) {
	error_occured = errors.New(msg)
	add_bit_state(bit, msg)
}

/*
   This Struct represents a square in the world.
   It has a color and a boolean that represents whether or not it has a bit.
*/
type Square struct {
	color Color
	has_bit bool
}

/*
   This function loads the world from a file.
   It returns a 2d array of squares.
   The format of the file is as follows:
   Use 'u' 'd' 'l' 'r' to represent the direction that the bit will start facing.
   then use '-' to represent a white square, 'r' to represent a red square, 'b' to represent a blue square, 'g' to represent a green square, and 'x' to represent a black square.
   Then specify the starting position of the bit with coordinates x y.

Example:
   u
   -r-r-
   -b-b-
   -r-r-
   0 0
   
*/
func load_world(file_name string) [][]Square {
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var world [][]Square
	var start_x int = -1
	var start_y int = -1
	for scanner.Scan() {
		var row []Square
		var broken bool = false
		for _, char := range scanner.Text() {
			if char == '-' {
				row = append(row, Square{color: White, has_bit: false})
			} else if char == 'r' {
				row = append(row, Square{color: Red, has_bit: false})
			} else if char == 'b' {
				row = append(row, Square{color: Blue, has_bit: false})
			} else if char == 'g' {
				row = append(row, Square{color: Green, has_bit: false})
			} else if char == 'x' {
				row = append(row, Square{color: Black, has_bit: false})
			} else if char == 'u' || char == 'd' || char == 'l' || char == 'r' {
				broken = true
				break
			} else if unicode.IsDigit(rune(char)) {
				if start_x == -1 {
					start_x = int(char) - 48
				} else {
					start_y = int(char) - 48
				}
			}
		}
		if !broken {
			world = append(world, row)
		}
	}
	world[start_y][start_x].has_bit = true
	world_height = len(world)
	world_width = len(world[0])
	return world
}


/*
   This function returns a bit and takes in the starting world and the ending world.
   The starting world is the world that bit starts in.
   The ending world is the world that bit must reach.
   This function should be called to get a bit.
*/
func GetBit(start_world string, end_world string) *Bit {
	reset_bit_globals()
	file, err := os.Open(start_world)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var bit_direction string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, char := range scanner.Text() {
			if char == 'u' {
				bit_direction = BitUp
			} else if char == 'd' {
				bit_direction = BitDown
			} else if char == 'l' {
				bit_direction = BitLeft
			} else if char == 'r' {
				bit_direction = BitRight
			} else {
				break
			}
		}
	}
	file.Close()

	var initial_world [][]Square = load_world(start_world)
	var final_world [][]Square = load_world(end_world)

	bit := new_bit(bit_direction, initial_world, final_world)
	init_bit_state(bit)
	return bit
}

/*
   Helper function for locating bit for the initial world.
*/
func find_bit_x(world [][]Square) int {
	for _, row := range world {
		for j, square := range row {
			if square.has_bit {
				return j
			}
		}
	}
	return 0
}

/*
   Helper function for locating bit for the initial world.
*/
func find_bit_y(world [][]Square) int {
	for i, row := range world {
		for _, square := range row {
			if square.has_bit {
				return i
			}
		}
	}
	return 0
}
		
	
/*
   This struct represents a bit.
   It has a face that represents the direction that it is facing.
   It has a number of steps that represents the number of steps it has taken.
   It has an x and y coordinate that represents its position in the world.
   It has a world that represents the world that it is currently in.
   It has a final_world that represents the world that it must reach.
*/
type Bit struct {
	face string
	steps int
	x int
	y int
	world [][]Square
	final_world [][]Square
}

/*
   This creates a new bit pointer.
   This function should not be called directly.
*/
func new_bit(direction string, world [][]Square, final_world[][]Square) *Bit {
	return &Bit{
		face: direction,
		steps: 0,
		x: find_bit_x(world),
		y: find_bit_y(world),
		world: world,
		final_world: final_world,
	}
}


/*
   This function prints a bit to the screen.
   It takes in a view and a bit face and a world.
*/ 
func print_world(v *gocui.View, face string, world [][]Square) {
	for _, row := range world {
			for i := 0; i < 3; i++ {
				for _, square := range row {
				
					if square.has_bit && i == 1 {
						switch square.color {
						case White:
							fmt.Fprint(v, "\x1b[36;47m  ",face, "  ")
						case Red:
							fmt.Fprint(v, "\x1b[36;41m  ",face, "  ")
						case Blue:
							fmt.Fprint(v, "\x1b[36;44m  ",face, "  ")
						case Green:
							fmt.Fprint(v, "\x1b[36;42m  ",face, "  ")
						default: 
							fmt.Fprint(v, "\x1b[36;47m  ",face, "  ")
						}
					} else {
						switch square.color {
						case White:
							fmt.Fprint(v, "\x1b[36;47m     ")
						case Black:
							fmt.Fprint(v, "\x1b[36;40m     ")
						case Red:
							fmt.Fprint(v, "\x1b[36;41m     ")
						case Blue:
							fmt.Fprint(v, "\x1b[36;44m     ")
						case Green:
							fmt.Fprint(v, "\x1b[36;42m     ")
						default: 
							fmt.Fprint(v, "\x1b[36;47m     ")
						}
					}
				}
			fmt.Fprintln(v)
		}
	}
}

/*
   This function moves bit in the direction that it is facing.
   It will mark an error if bit tries to move out of bounds or onto a black square.
*/
func (bit *Bit) Move() {
	if error_occured != nil {
		return
	}
	bit.steps++
	var result error
	switch bit.face {
	case BitUp:
		result = bit.moveUp()
	case BitDown:
		result = bit.moveDown()
	case BitLeft:
		result = bit.moveLeft()
	case BitRight:
		result = bit.moveRight()
	}
	if result != nil {
		stop_display_error(result.Error(),bit)
		return
	}
	add_bit_state(bit, "move     ")
}

/*
   Internal function for moving bit left.
   Returns an error if bit tries to make an invalid move.
*/
func (b *Bit) moveUp() error {
	if b.y -1 < 0 {
		return errors.New("Out of bounds")
	} else if b.world[b.y-1][b.x].color == White || b.world[b.y-1][b.x].color == Red || b.world[b.y-1][b.x].color == Blue || b.world[b.y-1][b.x].color == Green {
		b.world[b.y][b.x].has_bit = false
		b.world[b.y-1][b.x].has_bit = true
		b.y--
		return nil
	} else if b.world[b.y-1][b.x].color == Black {
		return errors.New("Blocked")
	}
	return nil
}

/*
   Internal function for moving bit down.
   Returns an error if bit tries to make an invalid move.
*/
func (b *Bit) moveDown() error {
	if b.y +1 > len(b.world) {
		return errors.New("Out of bounds")
	} else if b.world[b.y+1][b.x].color == White || b.world[b.y+1][b.x].color == Red || b.world[b.y+1][b.x].color == Blue || b.world[b.y+1][b.x].color == Green {
		b.world[b.y][b.x].has_bit = false
		b.world[b.y+1][b.x].has_bit = true
		b.y++
		return nil
	} else if b.world[b.y+1][b.x].color == Black {
		return errors.New("Blocked")
	}
	return nil
}

/*
   Internal function for moving bit left.
   Returns an error if bit tries to make an invalid move.
*/
func (b *Bit) moveLeft() error {
	if b.x -1 < 0 {
		return errors.New("Out of bounds")
	} else if b.world[b.y][b.x-1].color == White || b.world[b.y][b.x-1].color == Red || b.world[b.y][b.x-1].color == Blue || b.world[b.y][b.x-1].color == Green {
		b.world[b.y][b.x].has_bit = false
		b.world[b.y][b.x-1].has_bit = true
		b.x--
		return nil
	} else if b.world[b.y][b.x-1].color == Black {
		return errors.New("Blocked")
	}
	return nil
}

/*
   Internal function for moving bit right.
   Returns an error if bit tries to make an invalid move.
*/
func (b *Bit) moveRight() error {
	if b.x +1 > len(b.world[0]) -1 {
		return errors.New("Out of bounds")
	} else if b.world[b.y][b.x+1].color == White || b.world[b.y][b.x+1].color == Red || b.world[b.y][b.x+1].color == Blue || b.world[b.y][b.x+1].color == Green {
		b.world[b.y][b.x].has_bit = false
		b.world[b.y][b.x+1].has_bit = true
		b.x++
		return nil
	} else if b.world[b.y][b.x+1].color == Black {
		return errors.New("Blocked")
	}
	return nil
}

/*
   This function makes bit turn left without moving.
*/
func (b *Bit) Left() {
	if error_occured != nil {
		return
	}
	b.steps++
	switch b.face {
	case BitUp:
		b.face = BitLeft
	case BitDown:
		b.face = BitRight
	case BitLeft:
		b.face = BitDown
	case BitRight:
		b.face = BitUp
	}
	add_bit_state(b, "left       ")
	/*fmt.Println("Step  #", b.steps)
	b.print_world()
	pause()*/
}

/*
   This function makes bit turn right without moving.
*/
func (b *Bit) Right() {
	if error_occured != nil {
		return
	}
	b.steps++
	switch b.face {
	case BitUp:
		b.face = BitRight
	case BitDown:
		b.face = BitLeft
	case BitLeft:
		b.face = BitUp
	case BitRight:
		b.face = BitDown
	}
	add_bit_state(b, "right       ")
}

/*
   This function causes bit to paint a color on the current square with a given color.
   Valid colors are "red", "blue", and "green".
*/
func (bit *Bit) Paint(color string) {
	if error_occured != nil {
		return
	}
	var success bool = false
	if color == "red" {
		bit.world[bit.y][bit.x].color = Red
		success = true
	}
	if color == "blue" {
		bit.world[bit.y][bit.x].color = Blue
		success = true
	}
	if color == "green" {
		bit.world[bit.y][bit.x].color = Green
		success = true
	} else {
		stop_display_error("Invalid color: " + color,bit)
	}
	if success {
		add_bit_state(bit, "paint " + color)
	}
}

/*
   This function causes bit to erase the color on the current square.
*/
func (bit *Bit) Erase() {
	if error_occured != nil {
		return
	}
	bit.world[bit.y][bit.x].color = White
	add_bit_state(bit, "erase")
}

/*
   This function returns the color of the current square.
   Valid colors are "red", "blue", "green", and "white".
*/
func (bit *Bit) GetColor() string {
	if error_occured != nil {
		return ""
	}
	add_bit_state(bit, "get color       ")
	switch bit.world[bit.y][bit.x].color {
	case Red:
		return "red"
	case Blue:
		return "blue"
	case Green:
		return "green"
	default:
		return "white"
	}
}

/*
   This function checks if the the current square is red.
*/
func (bit *Bit) IsRed() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is red       ")
	return bit.world[bit.y][bit.x].color == Red
}

/*
   This function checks if the the current square is blue.
*/
func (bit *Bit) IsBlue() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is blue       ")
	return bit.world[bit.y][bit.x].color == Blue
}

/*
   This function checks if the the current square is green.
*/
func (bit *Bit) IsGreen() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is green       ")
	return bit.world[bit.y][bit.x].color == Green
}

/*
   This function checks if the square infront of bit is clear.
   Meaning that it is not black and is within the bounds of the world.
*/
func (bit *Bit) IsFrontClear() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is front clear")
	switch bit.face {
	case BitUp:
		return bit.world[bit.y-1][bit.x].color != Black || bit.y-1 < 0 
	case BitDown:
		return bit.world[bit.y+1][bit.x].color != Black || bit.y+1 >= len(bit.world)
	case BitLeft:
		return bit.world[bit.y][bit.x-1].color != Black || bit.x-1 < 0
	case BitRight:
		return bit.world[bit.y][bit.x+1].color != Black || bit.x+1 >= len(bit.world)
	}
	return false
}

/*
   This function checks if the square to the right of bit is clear.
   Meaning that it is not black and is within the bounds of the world.
*/
func (bit *Bit) IsRightClear() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is right clear")
	switch bit.face {
	case BitUp:
		return bit.world[bit.y][bit.x+1].color != Black || bit.x+1 >= len(bit.world)
	case BitDown:
		return bit.world[bit.y][bit.x-1].color != Black || bit.x-1 < 0
	case BitLeft:
		return bit.world[bit.y-1][bit.x].color != Black || bit.y-1 < 0
	case BitRight:
		return bit.world[bit.y+1][bit.x].color != Black || bit.y+1 >= len(bit.world)
	}
	return false
}

/*
   This function checks if the square to the left of bit is clear.
   Meaning that it is not black and is within the bounds of the world.
*/
func (bit *Bit) IsLeftClear() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is left clear")
	switch bit.face {
	case BitUp:
		return bit.world[bit.y][bit.x-1].color != Black || bit.x-1 < 0
	case BitDown:
		return bit.world[bit.y][bit.x+1].color != Black || bit.x-1 >= len(bit.world)
	case BitLeft:
		return bit.world[bit.y+1][bit.x].color != Black || bit.y+1 >= len(bit.world)
	case BitRight:
		return bit.world[bit.y-1][bit.x].color != Black || bit.y-1 < 0
	}
	return false
}

/*
   This function creates a new snapshot of the current state of the world.
   It takes a string as a parameter which is the name of the snapshot.
*/
func (bit *Bit) Snapshot(name string) {
	if error_occured != nil {
		return
	}
	add_bit_state(bit, "snapshot " + name)
	bit_snapshots = append(bit_snapshots, bit.steps)
	bit_snapshot_names = append(bit_snapshot_names, name)
}

/*
   This function checks to see if the current state of the world matches the final state of the world.
*/
func (bit *Bit) Compare() {
	var matches bool = true

	if error_occured != nil {
		return
	}
	
	for i := range bit.world {
		for j := range bit.world[i] {
			if bit.world[i][j].color != bit.final_world[i][j].color || bit.world[i][j].has_bit != bit.final_world[i][j].has_bit {
				matches = false
				break
			}
		}
	}
	if matches {
		add_bit_state(bit, "Success!   ")
	} else {
		add_bit_state(bit, "Error: Does not match")
	}
}
