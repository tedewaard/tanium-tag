package main

type Tag struct {
	name        string
	description string
}

//Examples of core tags. One of these should be on each PC
var coreTags = []Tag{
	{name: "User", description: "Computer is assigned to a specific user. Not a shared PC or critical PC."},
	{name: "Manufacturing", description: "Shared or generic Manufacturing PC."},
	{name: "Critical", description: "Critical PC that should not be managed. No patching, updates or software deployments."},
}

//This provides options on what kind of tag we want to apply to the machine
var typeOfTag = []Tag{
	{name: "Core PC Tagging", description: "Apply core tags like Business Unit and critical."},
	{name: "Application Tagging", description: "Apply application tags to install software."},
}

// This array can be updated with any tags we want to provide as options
var rbacTags = []Tag{
	{name: "rbac-exampleTag1", description: "Any computer that isn't a server and should be managed."},
	{name: "rbac-exampleTag2", description: "Any server"},
}

//Example Business Unit tags to identify where the PC belongs
var buTags = []Tag{
	{name: "BU-Marketing", description: ""},
	{name: "BU-IT", description: ""},
	{name: "BU-Engineering", description: ""},
	{name: "BU-ServiceDesk", description: ""},
	{name: "BU-HR", description: ""},
	{name: "BU-ELT", description: ""},
}

//Example application tags used to provide a PC access to a certain application
var appTags = []Tag{
	{name: "app-Chrome", description: ""},
	{name: "app-Firefox", description: ""},
	{name: "app-Edge", description: ""},
	{name: "app-Adobe", description: ""},
	{name: "app-WireShark", description: ""},
}
