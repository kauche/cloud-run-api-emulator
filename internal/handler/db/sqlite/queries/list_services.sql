SELECT
    parent,
    name,
    description,
    uid,
    generation,
    created_at
  FROM services
  WHERE
    parent = %%parent string%%
  ORDER BY
    created_at DESC,
    parent ASC,
    name ASC
  LIMIT
    %%limit int32%%;
