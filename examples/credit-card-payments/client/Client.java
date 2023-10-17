import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.util.Scanner;
import ar.com.montepagano.search.middleware.v1.PrivateMiddlewareServiceGrpc;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelRequest;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelResponse;

public class Client {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);

        // List of book titles
        String[] bookTitles = {
            "The Great Gatsby",
            "To Kill a Mockingbird",
            "1984",
            "Pride and Prejudice",
            "Animal Farm",
            "The Catcher in the Rye",
            "Lord of the Flies",
            "Brave New World",
            "The Hobbit",
            "The Fellowship of the Ring"
        };

        // Prompt user to select book titles
        System.out.println("Please select at least one book to purchase:");
        for (int i = 0; i < bookTitles.length; i++) {
            System.out.println((i + 1) + ". " + bookTitles[i]);
        }

        // Get user input for selected book titles
        System.out.print("Enter the number(s) of the book(s) you want to purchase (comma-separated): ");
        String selectedBooksInput = scanner.nextLine();
        String[] selectedBooksInputArray = selectedBooksInput.split(",");
        int[] selectedBooks = new int[selectedBooksInputArray.length];
        for (int i = 0; i < selectedBooksInputArray.length; i++) {
            selectedBooks[i] = Integer.parseInt(selectedBooksInputArray[i].trim());
        }

        // Prompt user for shipping address
        System.out.print("Enter your shipping address: ");
        String shippingAddress = scanner.nextLine();

        // Print purchase information
        System.out.println("Purchase Information:");
        System.out.println("Selected Books:");
        for (int i = 0; i < selectedBooks.length; i++) {
            System.out.println("- " + bookTitles[selectedBooks[i] - 1]);
        }
        System.out.println("Shipping Address: " + shippingAddress);

        // Register channel with middleware
        ManagedChannel channel = ManagedChannelBuilder.forTarget("localhost:1234").usePlaintext().build();
        PrivateMiddlewareServiceGrpc.PrivateMiddlewareServiceBlockingStub stub = PrivateMiddlewareServiceGrpc.newBlockingStub(channel);
        RegisterChannelRequest request = RegisterChannelRequest.newBuilder().build(); // Fill in the contract
        RegisterChannelResponse response = stub.registerChannel(request);
        channel.shutdown();
    }
}
