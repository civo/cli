package cmd

// volume list -- list all volumes [ls, all]
// volume create NAME SIZE -- create a volume of SIZE (GB) called NAME [new]
// volume resize ID NEW_SIZE -- resizes the volume with ID to NEW_SIZE GB
// volume remove ID -- remove the volume with ID [delete, destroy, rm]
// volume attach VOLUME_ID INSTANCE_ID -- connect the volume with VOLUME_ID to the instance with INSTANCE_ID [connect, link]
// volume detach ID -- disconnect the volume with ID from any instance it's connected to [disconnect, unlink]
// volume rename ID NEW_NAME -- rename the volume with ID to have the NEW_NAME
