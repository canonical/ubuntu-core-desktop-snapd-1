// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
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

import (
	"github.com/snapcore/snapd/osutil"
)

const polkitAgentSummary = `allows operation as a polkit agent`

const polkitAgentBaseDeclarationPlugs = `
  polkit-agent:
    allow-installation: false
    deny-auto-connection: true
`

const polkitAgentBaseDeclarationSlots = `
  polkit-agent:
    allow-installation:
      slot-snap-type:
        - core
    deny-auto-connection: true
`

const polkitAgentConnectedPlugAppArmor = `
# Description: Allow registering with polkitd as a polkit agent.

#include <abstractions/dbus-strict>

# Allow application to register as an agent with polkitd
dbus (receive, send)
    bus=system
    path=/org/freedesktop/PolicyKit1/Authority
    interface=org.freedesktop.PolicyKit1.Authority
    peer=(label=unconfined),
dbus (send)
    bus=system
    path=/org/freedesktop/PolicyKit1/Authority
    interface=org.freedesktop.DBus.Properties
    peer=(label=unconfined),
dbus (send)
    bus=system
    path=/org/freedesktop/PolicyKit1/Authority
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(label=unconfined),

# Allow polkitd to communicate with the agent
dbus (receive)
    bus=system
    interface=org.freedesktop.PolicyKit1.AuthenticationAgent
    peer=(label=unconfined),

# Allow communication with accounts-daemon. This is used by
# gnome-shell's agent implementation to display user information in
# the authorisation dialog.
dbus (send)
    bus=system
    path=/org/freedesktop/Accounts
    interface=org.freedesktop.DBus.Introspectable
    member=Introspect
    peer=(label=unconfined),

dbus (receive, send)
    bus=system
    path=/org/freedesktop/Accounts
    interface=org.freedesktop.Accounts
    member=FindUserByName
    peer=(label=unconfined),

dbus (receive, send)
    bus=system
    path=/org/freedesktop/Accounts/User[0-9]*
    interface=org.freedesktop.DBus.Properties
    member={Get,GetAll,PropertiesChanged}
    peer=(label=unconfined),

# Allow agent to execute the setuid polkit-agent-helper-1 unconfined
# TODO: determine whether this could run as a sub-profile
/usr/{libexec,lib/policykit-1}/polkit-agent-helper-1 Uxr,
`

const polkitAgentConnectedPlugSecComp = `
# Description: Allow polkit-agent-helper-1 to use the audit subsystem
bind
socket AF_NETLINK - NETLINK_AUDIT
`

func init() {
	registerIface(&commonInterface{
		name:                  "polkit-agent",
		summary:               polkitAgentSummary,
		implicitOnCore:        osutil.FileExists("/usr/libexec/polkit-agent-helper-1") || osutil.FileExists("/usr/lib/policykit-1/polkit-agent-helper-1"),
		baseDeclarationPlugs:  polkitAgentBaseDeclarationPlugs,
		baseDeclarationSlots:  polkitAgentBaseDeclarationSlots,
		connectedPlugAppArmor: polkitAgentConnectedPlugAppArmor,
		connectedPlugSecComp:  polkitAgentConnectedPlugSecComp,
	})
}
