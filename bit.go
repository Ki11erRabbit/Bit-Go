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

const WhiteSquare = "⬜"
const BlackSquare = "⬛"
const RedSquare = "🟥"
const BlueSquare = "🟦"
const GreenSquare = "🟩"

type Color int

const (
	White Color = iota
	Black
	Red
	Blue
	Green
)

func update_step_view() {
	gui_i.DeleteView("Steps")
	if step_v, err := gui_i.SetView("Steps", max_x/2-len(bit_state_actions[current_state]),max_y/2-7,max_x/2+len(bit_state_actions[current_state]), max_y/2-5); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		step_view = step_v
		fmt.Fprintln(step_view, current_state, ": ", bit_state_actions[current_state])
	}
}

func first_state(gui *gocui.Gui, world *gocui.View) error {
	world.Clear()
	current_state = 0
	print_world(world, bit_states[current_state].face, bit_states[current_state].world)
	update_step_view()
	return nil
}

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

func last_state(gui *gocui.Gui, world *gocui.View) error {
	world.Clear()
	current_state = len(bit_states) - 1
	print_world(world, bit_states[current_state].face, bit_states[current_state].world)
	update_step_view()
	return nil
}

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

var final_world bool = false

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

func quit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

var gui_i *gocui.Gui = nil
var max_x int = 0
var max_y int = 0

func RunGui() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer gui.Close()
	
	gui_i = gui

	max_x, max_y = gui.Size()

	var x0, y0, x1, y1 int
	
	if world_width % 2 != 0 {
		x0 = max_x/2 - (world_width/2)
		x1 = max_x/2 + (world_width/2) +1
	} else {
		x0 = max_x/2 - (world_width/2)
		x1 = max_x/2 + (world_width/2)
	}
	if world_height % 2 != 0 {
		y0 = max_y/2 - (world_height/2)
		y1 = max_y/2 + (world_height/2) +1
	} else {
		y0 = max_y/2 - (world_height/2)
		y1 = max_y/2 + (world_height/2)
	}
		
	
	if world, err := gui.SetView("World", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		current_state = len(bit_states) - 1

		print_world(world, bit_states[current_state].face, bit_states[current_state].world)

		if err := gui.SetKeybinding("World", rune('n'), gocui.ModNone, next_state); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := gui.SetKeybinding("World", rune('p'), gocui.ModNone, prev_state); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := gui.SetKeybinding("World", rune('f'), gocui.ModNone, first_state); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := gui.SetKeybinding("World", rune('l'), gocui.ModNone, last_state); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := gui.SetKeybinding("", rune('q'), gocui.ModNone, quit); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := gui.SetKeybinding("World", rune('s'), gocui.ModNone, switch_world); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if _, err := gui.SetCurrentView("World"); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
		}

		
	}


	if step_v, err := gui.SetView("Steps", max_x/2-len(bit_state_actions[current_state]),max_y/2-7,max_x/2+len(bit_state_actions[current_state]), max_y/2-5); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		step_view = step_v


		fmt.Fprintln(step_view, current_state, ": ", bit_state_actions[current_state])
	}

	if help_view, err := gui.SetView("Help", max_x/2 - len(HelpText)/2 -1, max_y/2+(world_height), max_x/2+len(HelpText)/2 +1, max_y/2+(world_height) +2); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Fprintln(help_view, "n: next step    p: previous step    f: first step    l: last step    s: switch world    q: quit")
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		os.Exit(1)
	}

}

var step_view *gocui.View = nil


var bit_snapshots []int = make([]int, 0)

var current_snapshot int = 0

var current_state int = 0
var bit_states []*Bit = make([]*Bit, 1)
var bit_state_actions []string = make([]string, 1)

var error_occured error = nil

var world_width int = 0
var world_height int = 0

func init_bit_state(bit *Bit) {
	dest := make([][]Square, len(bit.world))
	for i := range dest {
		dest[i] = make([]Square, len(bit.world[i]))
		copy(dest[i], bit.world[i])
	}
	bit_states[0] = &Bit{ face: bit.face, steps: bit.steps, x: bit.x, y: bit.y, world: dest, final_world: bit.final_world}
	bit_state_actions[0] = "initial state"
}

func add_bit_state(bit *Bit, action string) {
	dest := make([][]Square, len(bit.world))
	for i := range dest {
		dest[i] = make([]Square, len(bit.world[i]))
		copy(dest[i], bit.world[i])
	}
	bit_states = append(bit_states, &Bit{ face: bit.face, steps: bit.steps, x: bit.x, y: bit.y, world: dest, final_world: bit.final_world})
	bit_state_actions = append(bit_state_actions, action)
}

func stop_display_error(msg string, bit *Bit) {
	error_occured = errors.New(msg)
	add_bit_state(bit, msg)
}

type Square struct {
	color Color
	has_bit bool
}

var print_pair bool = false

func set_print_pair() {
	print_pair = true
}

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

func GetBit(start_world string, end_world string) *Bit {
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
		
	

type Bit struct {
	face string
	steps int
	x int
	y int
	world [][]Square
	final_world [][]Square
}

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

func print_world(v *gocui.View, face string, world [][]Square) {
	for _, row := range world {
		for _, square := range row {
			if square.has_bit {
				switch square.color {
				case White:
					fmt.Fprint(v, "\x1b[36;47m",face)
				case Red:
					fmt.Fprint(v, "\x1b[36;41m",face)
				case Blue:
					fmt.Fprint(v, "\x1b[36;44m",face)
				case Green:
					fmt.Fprint(v, "\x1b[36;42m",face)
				default: 
					fmt.Fprint(v, "\x1b[36;47m",face)
				}
			} else {
				switch square.color {
				case White:
					fmt.Fprint(v, "\x1b[36;47m ")
				case Black:
					fmt.Fprint(v, "\x1b[36;40m ")
				case Red:
					fmt.Fprint(v, "\x1b[36;41m ")
				case Blue:
					fmt.Fprint(v, "\x1b[36;44m ")
				case Green:
					fmt.Fprint(v, "\x1b[36;42m ")
				default: 
					fmt.Fprint(v, "\x1b[36;47m ")
				}
			}
		}
		fmt.Fprintln(v)
	}
}

/*func (bit *bit) print_world_with_bit() {
	for _, row := range bit.world {
		for _, square := range row {
			if square.has_bit {
				fmt.Print(bit.face)
			} else {
				fmt.Print(square.symbol)
			}
		}
		fmt.Println("\n")
	}
}

func (bit *bit) print_world_without_bit() {
	for _, row := range bit.world {
		for _, square := range row {
			fmt.Print(square.symbol)
		}
		fmt.Println("\n")
	}
}

func (bit *bit) print_world() {
	if bit.
	bit.print_world_with_bit()
	fmt.Println("--------------------")
	bit.print_world_without_bit()
}*/

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
	/*fmt.Println("Step  #", b.steps)
	b.print_world()
	if result != nil {
		stop_display_error(result.Error())
	}
	fmt.Println("--------------------")
	pause()*/
}

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
	/*fmt.Println("Step  #", b.steps)
	print_world(b.world)
	fmt.Println("--------------------")
	pause()*/
}

func (bit *Bit) Paint(color string) {
	if error_occured != nil {
		return
	}
	if color == "red" {
		bit.world[bit.y][bit.x].color = Red
	}
	if color == "blue" {
		bit.world[bit.y][bit.x].color = Blue
	}
	if color == "green" {
		bit.world[bit.y][bit.x].color = Green
	}
	add_bit_state(bit, "paint " + color)
}

func (bit *Bit) Erase() {
	if error_occured != nil {
		return
	}
	bit.world[bit.y][bit.x].color = White
	add_bit_state(bit, "erase")
}

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

func (bit *Bit) IsRed() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is red       ")
	return bit.world[bit.y][bit.x].color == Red
}

func (bit *Bit) IsRlue() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is blue       ")
	return bit.world[bit.y][bit.x].color == Blue
}

func (bit *Bit) IsGreen() bool {
	if error_occured != nil {
		return false
	}
	add_bit_state(bit, "is green       ")
	return bit.world[bit.y][bit.x].color == Green
}

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


func (bit *Bit) Snapshot() {
	if error_occured != nil {
		return
	}
	add_bit_state(bit, "snapshot       ")
	bit_snapshots = append(bit_snapshots, bit.steps)
}


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
