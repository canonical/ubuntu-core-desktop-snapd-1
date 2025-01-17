// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2019 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin

const desktopLaunchSummary = `allows snaps to identify and launch desktop applications in (or from) other snaps`

const desktopLaunchBaseDeclarationPlugs = `
  desktop-launch:
    allow-installation: false
    deny-auto-connection: true
`

const desktopLaunchBaseDeclarationSlots = `
  desktop-launch:
    allow-installation:
      slot-snap-type:
        - core
    deny-auto-connection: true
`

const desktopLaunchConnectedPlugAppArmor = `
# Description: Can identify and launch other snaps.

# Access to the desktop and icon files installed by snaps
/var/lib/snapd/desktop/applications/{,*} r,
/var/lib/snapd/desktop/icons/{,**} r,

# Allow to execute snap commands.
/snap/bin/* ixr,

# Allow to execute the snap command, which is used to launch each application
/usr/bin/snap ixr,
/snap/snapd/*/usr/bin/snap ixr,

# Allow access to all read-only information provided by snaps
/snap/*/*/** r,

#include <abstractions/dbus-session-strict>

dbus (send)
    bus=session
    path=/io/snapcraft/PrivilegedDesktopLauncher
    interface=io.snapcraft.PrivilegedDesktopLauncher
    member=OpenDesktopEntry
    peer=(label=unconfined),
`

// Only implicitOnClassic since userd isn't yet usable on core
func init() {
	registerIface(&commonInterface{
		name:                  "desktop-launch",
		summary:               desktopLaunchSummary,
		implicitOnClassic:     true,
		implicitOnCore:        true,
		baseDeclarationPlugs:  desktopLaunchBaseDeclarationPlugs,
		baseDeclarationSlots:  desktopLaunchBaseDeclarationSlots,
		connectedPlugAppArmor: desktopLaunchConnectedPlugAppArmor,
	})
}
