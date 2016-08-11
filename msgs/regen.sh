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

## List of types to skip
## These have a field Type which conflicts with the ros.Message#Type() interface
#skip[0]="shape_msgs/SolidPrimitive.msg"
#skip[1]="sensor_msgs/JoyFeedback.msg"
#skip[2]="visualization_msgs/ImageMarker.msg"
#skip[3]="visualization_msgs/InteractiveMarkerUpdate.msg"
#skip[4]="visualization_msgs/Marker.msg"

for dir in $msgs_dirs; do
  echo "package $dir ..."
  mkdir -p $dir
  for file in $(find $share_dir/$dir/msg/ -d 1 -name '*.msg'); do
		target=$dir/${file##*/}
    if [[ " ${skip[@]} " =~ " ${target} " ]]; then
      echo "  SKIPPED: ${target}"
    else
      cp $file $target
      ros-gen-go msg --package=$dir --in=$file --out=$target.go
    fi
	done
done
