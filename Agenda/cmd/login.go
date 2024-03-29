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

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login -u [username] -p [password] -e [email] -c [phone]",
	Short: "login a user",
	Long: "login a user, a username and a password required",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		Log.SetPrefix("[Cmd]   ")
		Log.Printf("login --username=%s --password=%s", username, password)
		//fmt.Printf("[Cmd]   login --username=%s --password=%s \n", username, password)

		if err := entity.LoginUser(username, password); err != nil {
			Log.SetPrefix("[Error] ")
			Log.Println(err)
			//fmt.Print("[Error] ")
			//fmt.Println(err)
		} else {
			Log.SetPrefix("[OK]    ")
			Log.Println("login successfully")
			//fmt.Println("[OK]    login successfully")
		}
		//fmt.Println("login called")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	loginCmd.Flags().StringP("username", "u", "", "username")
	loginCmd.Flags().StringP("password", "p", "", "password")
}
