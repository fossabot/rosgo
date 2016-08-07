#!/bin/bash -e

share_dir="${ROS_ROOT}/.."
msgs_dirs='
	actionlib_msgs
	common_msgs
	control_msgs
	diagnostic_msgs
	geometry_msgs
	map_msgs
	nav_msgs
	rosgraph_msgs
	sensor_msgs
	shape_msgs
	smach_msgs
	std_msgs
	stereo_msgs
	tf2_msgs
	trajectory_msgs
	visualization_msgs
'

for dir in $msgs_dirs; do
  echo "package $dir ..."
  mkdir -p $dir
  for file in $(find $share_dir/$dir/msg/ -d 1 -name '*.msg'); do
		target=$dir/${file##*/}
		cp $file $target
		ros-gen-go msg --package=$dir --in=$file --out=$target.go
	done
done
