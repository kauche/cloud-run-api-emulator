SELECT
    parent,
    name,
    description,
    uid,
    uri,
    generation,
    created_at
  FROM services
  WHERE
    parent = %%parent string%% AND
    created_at <= %%created_at string%% AND
    (created_at < %%created_at string%% OR name > %%name string%%)
  ORDER BY
    created_at DESC,
    parent ASC,
    name ASC
  LIMIT
    %%limit int32%%;
