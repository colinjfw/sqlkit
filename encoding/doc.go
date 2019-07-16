// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

/*
Package encoding handles taking golang interfaces and converting them from
sql types to golang types. It uses sqlx under the hood to manage this.

One of the main differences is that it has a different philosophy from the core
golang database/sql package when it comes to nullable types. Mainly, nullable
types are converted into their default value in golang. Instead of requiring a
pointer.

You can call Marshal and Unmarshal directly or instantiate an Encoder. This is
a similar api to encoding/json and other golang packages.
*/
package encoding
