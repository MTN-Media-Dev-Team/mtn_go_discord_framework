package mtn_go_discord_framework

import "github.com/bwmarrin/discordgo"

type CommandOption interface {
	GetValue() interface{}
	GetName() string
}

type StringOption struct {
	Value string
	Name  string
}

func (s StringOption) GetValue() interface{} {
	return s.Value
}
func (s StringOption) GetName() string {
	return s.Name
}

type IntegerOption struct {
	Value int64
	Name  string
}

func (s IntegerOption) GetValue() interface{} {
	return s.Value
}

func (s IntegerOption) GetName() string {
	return s.Name
}

type UnsignedIntergerOption struct {
	Value uint64
	Name  string
}

func (s UnsignedIntergerOption) GetValue() interface{} {
	return s.Value
}

func (s UnsignedIntergerOption) GetName() string {
	return s.Name
}

type BooleanOption struct {
	Value bool
	Name  string
}

func (s BooleanOption) GetValue() interface{} {
	return s.Value
}

func (s BooleanOption) GetName() string {
	return s.Name
}

type FloatOption struct {
	Value float64
	Name  string
}

func (s FloatOption) GetValue() interface{} {
	return s.Value
}

func (s FloatOption) GetName() string {
	return s.Name
}

type UserOption struct {
	Value *discordgo.User
	Name  string
}

func (s UserOption) GetValue() interface{} {
	return s.Value
}

func (s UserOption) GetName() string {
	return s.Name
}

type ChannelOption struct {
	Value *discordgo.Channel
	Name  string
}

func (s ChannelOption) GetValue() interface{} {
	return s.Value
}

func (s ChannelOption) GetName() string {
	return s.Name
}

type RoleOption struct {
	Value *discordgo.Role
	Name  string
}

func (s RoleOption) GetValue() interface{} {
	return s.Value
}

func (s RoleOption) GetName() string {
	return s.Name
}
