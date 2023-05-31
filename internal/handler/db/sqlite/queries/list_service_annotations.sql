SELECT
    service_parent,
    service_name,
    key,
    value
  FROM service_annotations
  WHERE
    service_parent = %%parent string%% AND
    service_name = %%name string%%
