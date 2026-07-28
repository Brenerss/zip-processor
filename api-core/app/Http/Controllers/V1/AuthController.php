<?php

declare(strict_types = 1);

namespace App\Http\Controllers\V1;

use App\Http\Controllers\Controller;
use App\Models\User;
use Illuminate\Support\Facades\Auth;

class AuthController extends Controller
{
    public function loginTestUser()
    {
        /** @var User $user */
        $user = User::first();

        if (!$user) {
            return response()->json(['message' => 'No users found'], 404);
        }

        Auth::loginUsingId($user->id);

        return response()->json([
            'message' => 'Logged in as test user',
            'token'   => $user->createToken('test-user-token')->plainTextToken,
        ]);
    }
}
