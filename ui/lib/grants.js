export function sortByPrivilege(a, b) {
  if (a?.privilege === 'cluster-admin') {
    return -1
  }

  if (b?.privilege === 'cluster-admin') {
    return 1
  }

  return a?.privilege?.localeCompare(b?.privilege)
}

export function sortByResource(a, b) {
  return a?.resource?.localeCompare(b?.resource)
}

export function sortBySubject(a, b) {
  return (a?.user || a?.group)?.localeCompare(b?.user || b?.group)
}
