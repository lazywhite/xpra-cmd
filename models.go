//
// models.go
// Copyright (C) 2020 white <white@Whites-Mac-Air.local>
//
// Distributed under terms of the MIT license.
//

package main

type CommandInfo struct {
    CMD string `json:cmd`
}

type AllowedCommand struct{
    ID string `json:id`
    CreatedTS string `json:timestamp_created`
    UpdatedTS string `json:timestamp_updated`
    Label string `json:label`
    Command string `json:command`
    Category string `json:category`
    Icon string `json:icon`
    Rank int `json:rank`
    Enabled bool `json:enabled`
}
