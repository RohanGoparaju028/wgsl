### wgsl 
The whole idea of this project to build a command line tool with a set of commands that helps to do the predict and visualize leukimia cells present in the white blood cells and build in support to send the results to  email of the  doctor who are using it all via terminal or command prompt securely.
The list of supported commands are<br> 
<b> init </b> - initialise an empty wgsl that establishes an SSH connection to the server of the hospital/clinic and also verifies the email belongs to an actual person having medical experience.

<b> get </b> - This command is used to get securely the image that we wanna test using the machine learning to visualize and provide concrete evidence for the doctor to make concrete decision.

<b> train </b> - This command triggers the machine learning model for the images.

<b> test </b> - This command will test the new image.

<b> result </b> - Command to post the results to email 

<b> help </b> - Command to display all the commands.


