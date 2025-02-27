import Head from 'next/head'
import useSWR from 'swr'
import { useState, useRef } from 'react'
import dayjs from 'dayjs'
import { PlusIcon } from '@heroicons/react/outline'

import { useAdmin } from '../../lib/admin'

import Dashboard from '../../components/layouts/dashboard'
import PageHeader from '../../components/page-header'
import EmptyTable from '../../components/empty-table'
import Table from '../../components/table'
import Sidebar from '../../components/sidebar'
import EmptyData from '../../components/empty-data'
import IdentityList from '../../components/identity-list'
import TypeaheadCombobox from '../../components/typeahead-combobox'
import Metadata from '../../components/metadata'
import GrantsList from '../../components/grants-list'
import RemoveButton from '../../components/remove-button'

const columns = [
  {
    Header: 'Name',
    accessor: g => g,
    width: '67%',
    Cell: ({ value: group }) => {
      return (
        <div className='flex items-center py-1.5'>
          <div className='flex h-7 w-7 select-none items-center justify-center rounded-md border border-violet-300/40'>
            <img alt='group icon' src='/groups.svg' className='h-3 w-3' />
          </div>
          <div className='ml-3 flex min-w-0 flex-1 flex-col leading-tight'>
            <div className='truncate'>{group.name}</div>
          </div>
        </div>
      )
    },
  },
  {
    Header: 'Users',
    accessor: g => g,
    width: '33%',
    Cell: ({ value: group }) => {
      const { data: { items: users } = {}, error } = useSWR(
        `/api/users?group=${group.id}`
      )

      return (
        <>
          {users && (
            <div className='text-gray-400'>
              {users?.length} {users?.length === 1 ? 'member' : 'members'}
            </div>
          )}
          {error?.status && <div className='text-gray-400'>--</div>}
        </>
      )
    },
  },
]

function EmailsSelectInput({
  selectedEmails,
  setSelectedEmails,
  existMembers,
  onClick,
}) {
  const { data: { items: users } = { items: [] } } = useSWR('/api/users')

  const [query, setQuery] = useState('')
  const inputRef = useRef(null)

  const selectedEmailsId = selectedEmails.map(i => i.id)

  const filteredEmail = [...users.map(u => ({ ...u, user: true }))]
    .filter(s => s?.name?.toLowerCase()?.includes(query.toLowerCase()))
    .filter(s => s.name !== 'connector')
    .filter(s => !selectedEmailsId?.includes(s.id))
    .filter(s => !existMembers?.includes(s.id))

  const removeSelectedEmail = email => {
    setSelectedEmails(selectedEmails.filter(item => item.id !== email.id))
  }

  return (
    <section className='my-2 flex'>
      <div className='flex flex-1 items-center border-b border-gray-800 py-3'>
        <TypeaheadCombobox
          selectedEmails={selectedEmails}
          setSelectedEmails={setSelectedEmails}
          onRemove={removedEmail => removeSelectedEmail(removedEmail)}
          inputRef={inputRef}
          setQuery={setQuery}
          filteredEmail={filteredEmail}
          onKeyDownEvent={key => {
            if (key === 'Backspace' && inputRef.current.value.length === 0) {
              removeSelectedEmail(selectedEmails[selectedEmails.length - 1])
            }
          }}
        />
      </div>
      <div className='relative mt-3'>
        <button
          type='button'
          onClick={onClick}
          disabled={selectedEmails.length === 0}
          className='flex h-8 cursor-pointer items-center rounded-md border border-violet-300 px-3 py-3 text-2xs disabled:transform-none disabled:cursor-default disabled:opacity-30 disabled:transition-none sm:ml-4 sm:mt-0'
        >
          <PlusIcon className='mr-1.5 h-3 w-3' />
          <div className='text-violet-100'>Add</div>
        </button>
      </div>
    </section>
  )
}

function Details({ group, admin, onDelete }) {
  const { id, name } = group

  const { data: { items: users } = {}, mutate: mutateUsers } = useSWR(
    `/api/users?group=${group.id}`
  )
  const { data: { items } = {}, mutate: mutateGrants } = useSWR(
    `/api/grants?group=${id}`
  )

  const [emails, setEmails] = useState([])

  const grants = items?.filter(g => g.resource !== 'infra')
  const existMembers = users?.map(m => m.id)
  const metadata = [
    { title: 'ID', data: id },
    {
      title: 'Created',
      data: group?.created ? dayjs(group.created).fromNow() : '-',
    },
  ]

  return (
    <div className='flex flex-1 flex-col space-y-6'>
      {admin && (
        <>
          <section>
            <h3 className='mb-4 border-b border-gray-800 py-4 text-3xs uppercase text-gray-400'>
              Access
            </h3>
            <GrantsList
              grants={grants}
              onRemove={async id => {
                await fetch(`/api/grants/${id}`, { method: 'DELETE' })
                mutateGrants({ items: grants.filter(x => x.id !== id) })
              }}
              onChange={async (privilege, grant) => {
                const res = await fetch('/api/grants', {
                  method: 'POST',
                  body: JSON.stringify({
                    ...grant,
                    privilege,
                  }),
                })

                // delete old grant
                await fetch(`/api/grants/${grant.id}`, { method: 'DELETE' })
                mutateGrants({
                  items: [
                    ...grants.filter(f => f.id !== grant.id),
                    await res.json(),
                  ],
                })
              }}
            />
            {!grants?.length && (
              <EmptyData>
                <div className='mt-6'>No access</div>
              </EmptyData>
            )}
          </section>
          <section>
            <h3 className='mb-2 border-b border-gray-800 py-4 text-3xs uppercase text-gray-400'>
              Users
            </h3>
            <EmailsSelectInput
              selectedEmails={emails}
              setSelectedEmails={setEmails}
              existMembers={existMembers}
              onClick={async () => {
                const usersToAdd = emails.map(email => email.id)
                await fetch(`/api/groups/${id}/users`, {
                  method: 'PATCH',
                  body: JSON.stringify({ usersToAdd }),
                })

                await mutateUsers(`/api/users?group=${group.id}`)
                setEmails([])
              }}
            />
            <div className='mt-4'>
              {users?.length === 0 && (
                <EmptyData>
                  <div className='mt-6'>No members in the group</div>
                </EmptyData>
              )}
              <IdentityList
                list={users?.sort((a, b) =>
                  b.created?.localeCompare(a.created)
                )}
                actionText='Remove'
                onClick={async userId => {
                  const usersToRemove = [userId]
                  await fetch(`/api/groups/${id}/users`, {
                    method: 'PATCH',
                    body: JSON.stringify({ usersToRemove }),
                  })
                  mutateUsers({
                    items: users.filter(i => i.id !== userId),
                  })
                }}
              />
            </div>
          </section>
        </>
      )}
      <section>
        <h3 className='mb-4 border-b border-gray-800 py-4 text-3xs uppercase text-gray-400'>
          Metadata
        </h3>
        <Metadata data={metadata} />
      </section>
      {admin && (
        <section className='flex flex-1 flex-col items-end justify-end py-6'>
          <RemoveButton
            onRemove={async () => {
              onDelete()
            }}
            modalTitle='Remove Group'
            modalMessage={
              <>
                Are you sure you want to delete{' '}
                <span className='font-bold text-white'>{name}</span>? This
                action cannot be undone.
              </>
            }
          />
        </section>
      )}
    </div>
  )
}

export default function Groups() {
  const { data: { items: groups } = {}, error, mutate } = useSWR('/api/groups')
  const { admin, loading: adminLoading } = useAdmin()

  const [selected, setSelected] = useState(null)

  const loading = adminLoading || (!groups && !error)

  return (
    <>
      <Head>Groups - Infra</Head>
      {!loading && (
        <div className='flex h-full flex-1'>
          <div className='flex flex-1 flex-col space-y-4'>
            <PageHeader
              header='Groups'
              buttonHref='/groups/add'
              buttonLabel='Group'
            />
            {error?.status ? (
              <div className='my-20 text-center text-sm font-light text-gray-300'>
                {error?.info?.message}
              </div>
            ) : (
              <div className='flex min-h-0 flex-1 flex-col overflow-y-auto px-6'>
                <Table
                  columns={columns}
                  data={
                    groups?.sort((a, b) =>
                      b.created?.localeCompare(a.created)
                    ) || []
                  }
                  getRowProps={row => ({
                    onClick: () => setSelected(row.original),
                    className:
                      selected?.id === row.original.id
                        ? 'bg-gray-900/50'
                        : 'cursor-pointer',
                  })}
                />
                {groups?.length === 0 && (
                  <EmptyTable
                    title='There are no groups'
                    subtitle='Connect, create and manage your groups.'
                    iconPath='/groups.svg'
                    buttonHref='/groups/add'
                    buttonText='Groups'
                  />
                )}
              </div>
            )}
          </div>
          {selected && (
            <Sidebar
              onClose={() => setSelected(null)}
              title={selected?.name}
              iconPath='/groups.svg'
            >
              <Details
                group={selected}
                admin={admin}
                onDelete={() => {
                  mutate(async ({ items: groups } = { items: [] }) => {
                    await fetch(`/api/groups/${selected.id}`, {
                      method: 'DELETE',
                    })

                    return {
                      items: groups?.filter(g => g?.id !== selected.id),
                    }
                  })
                  setSelected(null)
                }}
              />
            </Sidebar>
          )}
        </div>
      )}
    </>
  )
}

Groups.layout = function (page) {
  return <Dashboard>{page}</Dashboard>
}
