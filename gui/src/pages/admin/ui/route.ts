import RouteViewComponent from '../../../layouts/RouterBypass.vue'

export default {
  name: 'ui',
  path: 'ui',
  component: RouteViewComponent,
  children: [
    {
      name: 'typography',
      path: 'typography',
      component: () => import('../../../pages/admin/ui/typography/Typography.vue'),
    },
  ],
}
