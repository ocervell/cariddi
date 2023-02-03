/*
==========
Cariddi
==========

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.

	@Repository:  https://github.com/ocervell/cariddi

	@Author:      edoardottt, https://www.edoardoottavianelli.it

	@License: https://github.com/ocervell/cariddi/blob/main/LICENSE

*/

package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Beautify prints the banner + version.
func Beautify() {
	banner1 := "                 _     _     _ _ \n"
	banner2 := "   ___ __ _ _ __(_) __| | __| (_)\n"
	banner3 := "  / __/ _` | '__| |/ _` |/ _` | |\n"
	banner4 := " | (_| (_| | |  | | (_| | (_| | |\n"
	banner5 := "  \\___\\__,_|_|  |_|\\__,_|\\__,_|_| v1.2.1\n"
	banner6 := "\n"
	banner7 := " > github.com/ocervell/cariddi\n"
	banner8 := " > edoardoottavianelli.it\n"
	banner9 := "========================================\n"

	bannerPart1 := banner1 + banner2 + banner3 + banner4 + banner5
	bannerPart2 := banner6 + banner7 + banner8 + banner9

	color.Set(color.FgCyan)
	fmt.Fprintf(os.Stderr, bannerPart1)
	color.Unset()
	fmt.Fprintf(os.Stderr, bannerPart2)
}
