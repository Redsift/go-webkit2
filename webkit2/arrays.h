#ifndef GO_ARRAYS_H
#define GO_ARRAYS_H

#include <stdlib.h>

static gchar** alloc_gchar_array(size_t l)
{
	gchar**	v;

	v = calloc(l, sizeof(gchar*));
	return v;
}

static void free_gchar_array(gchar** v)
{
	int	i;

	if (v == NULL) {
		return;
	}

	for (i = 0; v[i] != NULL; ++i) {
		free(v[i]);
	}
	free(v);
}

static void set_gchar_array(gchar** v, int i, gchar* s)
{
	v[i] = s;
}

#endif