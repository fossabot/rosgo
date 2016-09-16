#!/bin/bash -e

share_dir="${ROS_ROOT}/.."
srvs_dirs='
	std_srvs
'

for dir in $srvs_dirs; do
  echo "package $dir ..."
  mkdir -p $dir
  for file in $(find $share_dir/$dir/srv/ -d 1 -name '*.srv'); do
		target=$dir/${file##*/}
    cp $file $target
    ros-gen-go srv --package=$dir --in=$file --out=$target.go
	done
done
