echo "Test 3 - Motor A forward at 100%"

# Dir 1 Pin 13 = BCM 27
# Dir 2 Pin 15 = BCM 22
# Speed Pin 11 = BCM 17

# off
echo "17=0; 27=0; 22=0" > /dev/pi-blaster
sleep 2

echo "27=0; 22=1; 17=1" > /dev/pi-blaster
sleep 5

# off
echo "17=0; 27=0; 22=0" > /dev/pi-blaster
