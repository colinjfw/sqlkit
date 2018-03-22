#!/bin/sh

license="// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
"

for file in $(find ./db ./encoding | grep 'go'); do
  if [ "$(head -n 1 $file | grep 'Copyright')" = "" ]; then
    mv $file $file.bak
    echo "$license" > $file
    cat $file.bak >> $file
    rm $file.bak
  fi
done
