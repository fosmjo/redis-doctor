/*
Copyright © 2023 fosmjo <imefangjie@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"

	"github.com/fosmjo/redis-doctor/pkg/doctor"
	"github.com/fosmjo/redis-doctor/pkg/visitors"
)

type options struct {
	host     string
	port     int
	db       int
	user     string
	password string

	symptom string
	pattern string
	format  string
	count   int
}

func (o *options) toRedisUniversalOPtions() *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:      []string{o.host + ":" + strconv.Itoa(o.port)},
		ClientName: "redis-doctor",
		DB:         o.db,
		Username:   o.user,
		Password:   o.password,
	}
}

var _options = &options{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redis-doctor",
	Short: "diagnose redis problems.",
	Long:  `redis-doctor is a cli tool for diagnosing redis problems, such as hotkey, bigkey, slowlog, etc.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var outputer doctor.Visitor
		switch _options.format {
		case "json":
			outputer = visitors.NewJSONVisitor(os.Stdout)
		case "xml":
			v := visitors.NewXMLVisitor(os.Stdout)
			defer func() {
				err = v.Close()
			}()
			outputer = v
		default: // csv
			v := visitors.NewCSVVisitor(os.Stdout)
			defer v.Flush()
			outputer = v
		}

		d := doctor.New(_options.toRedisUniversalOPtions(), outputer)

		return d.Diagnose(
			context.Background(),
			_options.symptom,
			doctor.Options{
				Pattern: _options.pattern,
				Count:   _options.count,
			},
		)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&_options.host, "host", "", "127.0.0.1", "redis server host")
	rootCmd.Flags().IntVarP(&_options.port, "port", "p", 6379, "redis server port")
	rootCmd.Flags().IntVarP(&_options.db, "db", "n", 0, "redis database (default 0)")
	rootCmd.Flags().StringVarP(&_options.user, "user", "u", "", "redis username")
	rootCmd.Flags().StringVarP(&_options.password, "pass", "", "", "redis password")
	rootCmd.Flags().StringVarP(
		&_options.symptom, "symptom", "s", "",
		"symptom to diagnose (required, oneof: bigkey, hotkey, slowlog)",
	)
	err := rootCmd.MarkFlagRequired("symptom")
	if err != nil {
		panic(err)
	}
	rootCmd.Flags().StringVarP(
		&_options.pattern, "pattern", "", "*",
		"keys pattern when using the --bigkeys or --hotkey options",
	)
	rootCmd.Flags().IntVarP(&_options.count, "count", "c", 10, "specify the number of returned entries")
	rootCmd.Flags().StringVarP(
		&_options.format, "format", "f", "csv",
		"output format (oneof: csv, json, xml)",
	)
}
