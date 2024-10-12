#!/bin/bash

ffmpeg -f image2 -framerate 60 -i app/simulator/images/%05d.png -vcodec libx264 -crf 22 -y app/simulator/video.mp4
