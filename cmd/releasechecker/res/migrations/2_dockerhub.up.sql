CREATE TABLE dockerhub_images (
  dh_group TEXT,
  dh_name TEXT,
  PRIMARY KEY (dh_group, dh_name)
);

CREATE TABLE dockerhub_image_tags (
  dh_group TEXT,
  dh_name TEXT,
  tag_name TEXT,
  FOREIGN KEY (dh_group, dh_name) REFERENCES dockerhub_images(dh_group, dh_name),
  PRIMARY KEY (dh_group, dh_name, tag_name)
);