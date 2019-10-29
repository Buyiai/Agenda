/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	//"fmt"
	"github.com/github-user/Agenda/entity"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register -u [username] -p [password] -e [emial] -c [phone]",
	Short: "register a new user",
	Long: "register a new user, a unique username, a password, an email and a phone required",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")
		contact, _ := cmd.Flags().GetString("contact")

		Log.SetPrefix("[Cmd]   ")
		Log.Printf("register --username=%s --password=%s --email=%v --phone=%v", username, password, email, contact)
		//fmt.Printf("[Cmd]   register --username=%s --password=%s --email=%v --phone=%v \n", username, password, email, contact)

		if err := entity.RegisterUser(username, password, email, contact); err != nil {
			Log.SetPrefix("[Error] ")
			Log.Println(err)
			//fmt.Printf("[Error] ")
			//fmt.Println(err)
		} else {
			Log.SetPrefix("[OK]    ")
			Log.Println("register successfully")
			//fmt.Println("[OK]    register successfully")
		}
		//fmt.Println("register called")
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	registerCmd.Flags().StringP("username", "u", "", "username")
	registerCmd.Flags().StringP("password", "p", "", "password")
	registerCmd.Flags().StringP("email", "e", "", "email")
	registerCmd.Flags().StringP("contact", "c", "", "phone")
}
